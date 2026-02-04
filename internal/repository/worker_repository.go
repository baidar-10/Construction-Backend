package repository

import (
	"construction-backend/internal/database"
	"construction-backend/internal/models"

	"github.com/google/uuid"
)

type WorkerRepository struct {
	db *database.Database
}

func NewWorkerRepository(db *database.Database) *WorkerRepository {
	return &WorkerRepository{db: db}
}

func (r *WorkerRepository) Create(worker *models.Worker) error {
	return r.db.Create(worker).Error
}

func (r *WorkerRepository) FindAll(filters map[string]interface{}) ([]models.Worker, error) {
	var workers []models.Worker
	query := r.db.Preload("User").Preload("Portfolio").Preload("TeamMembers")

	if specialty, ok := filters["specialty"].(string); ok && specialty != "" {
		query = query.Where("specialty = ?", specialty)
	}
	if location, ok := filters["location"].(string); ok && location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}
	if minRate, ok := filters["minRate"].(float64); ok {
		query = query.Where("hourly_rate >= ?", minRate)
	}
	if maxRate, ok := filters["maxRate"].(float64); ok {
		query = query.Where("hourly_rate <= ?", maxRate)
	}
	if availability, ok := filters["availability"].(string); ok && availability != "" {
		query = query.Where("availability_status = ?", availability)
	}

	err := query.Find(&workers).Error
	if err != nil {
		return nil, err
	}

	// Load skills for each worker
	for i := range workers {
		var skills []models.WorkerSkill
		r.db.Where("worker_id = ?", workers[i].ID).Find(&skills)
		workers[i].Skills = make([]string, len(skills))
		for j, skill := range skills {
			workers[i].Skills[j] = skill.Skill
		}
	}

	return workers, nil
}

func (r *WorkerRepository) FindByID(id uuid.UUID) (*models.Worker, error) {
	var worker models.Worker
	err := r.db.Preload("User").Preload("Portfolio").Preload("TeamMembers").First(&worker, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	// Load skills
	var skills []models.WorkerSkill
	r.db.Where("worker_id = ?", worker.ID).Find(&skills)
	worker.Skills = make([]string, len(skills))
	for i, skill := range skills {
		worker.Skills[i] = skill.Skill
	}

	return &worker, nil
}

func (r *WorkerRepository) FindByUserID(userID uuid.UUID) (*models.Worker, error) {
	var worker models.Worker
	err := r.db.Preload("User").Preload("TeamMembers").Where("user_id = ?", userID).First(&worker).Error
	if err != nil {
		return &worker, err
	}

	// Load skills
	var skills []models.WorkerSkill
	r.db.Where("worker_id = ?", worker.ID).Find(&skills)
	worker.Skills = make([]string, len(skills))
	for i, skill := range skills {
		worker.Skills[i] = skill.Skill
	}

	return &worker, nil
}

func (r *WorkerRepository) Update(worker *models.Worker) error {
	return r.db.Save(worker).Error
}

func (r *WorkerRepository) Search(query string) ([]models.Worker, error) {
	var workers []models.Worker
	searchPattern := "%" + query + "%"
	
	err := r.db.Preload("User").Preload("Portfolio").Preload("TeamMembers").
		Joins("JOIN users ON users.id = workers.user_id").
		Where("users.first_name ILIKE ? OR users.last_name ILIKE ? OR workers.specialty ILIKE ? OR workers.location ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Find(&workers).Error
	
	if err != nil {
		return nil, err
	}

	// Load skills for each worker
	for i := range workers {
		var skills []models.WorkerSkill
		r.db.Where("worker_id = ?", workers[i].ID).Find(&skills)
		workers[i].Skills = make([]string, len(skills))
		for j, skill := range skills {
			workers[i].Skills[j] = skill.Skill
		}
	}

	return workers, nil
}

func (r *WorkerRepository) FilterBySkill(skill string) ([]models.Worker, error) {
	var workerSkills []models.WorkerSkill
	err := r.db.Where("skill = ?", skill).Find(&workerSkills).Error
	if err != nil {
		return nil, err
	}

	var workerIDs []uuid.UUID
	for _, ws := range workerSkills {
		workerIDs = append(workerIDs, ws.WorkerID)
	}

	var workers []models.Worker
	if len(workerIDs) > 0 {
			err = r.db.Preload("User").Preload("Portfolio").Preload("TeamMembers").Where("id IN ?", workerIDs).Find(&workers).Error
		if err != nil {
			return nil, err
		}

		// Load skills for each worker
		for i := range workers {
			var skills []models.WorkerSkill
			r.db.Where("worker_id = ?", workers[i].ID).Find(&skills)
			workers[i].Skills = make([]string, len(skills))
			for j, skill := range skills {
				workers[i].Skills[j] = skill.Skill
			}
		}
	}

	return workers, nil
}

func (r *WorkerRepository) AddSkill(workerID uuid.UUID, skill string) error {
	workerSkill := &models.WorkerSkill{
		WorkerID: workerID,
		Skill:    skill,
	}
	return r.db.Create(workerSkill).Error
}

func (r *WorkerRepository) AddPortfolio(portfolio *models.Portfolio) error {
	return r.db.Create(portfolio).Error
}