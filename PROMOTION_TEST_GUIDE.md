# Promotion System Testing Guide

## ✅ Алгоритм продвижения работает!

### API Endpoints

#### 1. **Просмотр прайса на продвижение** (Public)
```bash
curl http://localhost:8080/api/promotions/pricing
```

**Ответ:**
- `featured` - 100₸/день (Featured in search results)
- `top` - 250₸/день (Top position in search results)  
- `premium` - 500₸/день (Premium position with badge)

#### 2. **Просмотр ТОП воркеров** (Public)
```bash
curl http://localhost:8080/api/promotions/top-workers?limit=10
```

#### 3. **Продвинуть воркера** (Admin only)
```bash
curl -X POST http://localhost:8080/api/admin/workers/{workerId}/promote \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "promotionType": "premium",
    "durationDays": 30
  }'
```

#### 4. **Отменить продвижение** (Admin only)
```bash
curl -X POST http://localhost:8080/api/admin/workers/{workerId}/cancel-promotion \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

#### 5. **История продвижений воркера** (Public)
```bash
curl http://localhost:8080/api/promotions/workers/{workerId}/history
```

---

## 🎨 Как тестировать через UI

### 1. Войти как администратор
- Email: `admin@stroymaster.com`
- Password: `admin123`
- URL: http://localhost:5173/login

### 2. Зайти в админ-панель
- URL: http://localhost:5173/admin
- Раздел "Workers" или "Users"

### 3. Продвинуть воркера
На странице воркера должна быть кнопка/форма:
- **Выбрать тип продвижения:** featured / top / premium
- **Указать длительность:** 7-30 дней
- **Нажать "Promote"**

### 4. Проверить ТОП воркеров
- На главной странице `/workers` или `/`
- Должны отображаться сначала продвинутые воркеры
- С бейджами "PREMIUM" / "TOP" / "FEATURED"

---

## 📊 Проверка в БД

```bash
# Посмотреть всех продвинутых воркеров
docker exec construction_db psql -U admin -d construction_db -c "
  SELECT w.id, u.first_name, u.last_name, ph.promotion_type, ph.expires_at
  FROM workers w
  JOIN users u ON w.user_id = u.id
  LEFT JOIN promotion_history ph ON ph.worker_id::uuid = w.id AND ph.status = 'active'
  WHERE ph.id IS NOT NULL;
"

# Посмотреть историю продвижений
docker exec construction_db psql -U admin -d construction_db -c "
  SELECT * FROM promotion_history ORDER BY created_at DESC LIMIT 10;
"
```

---

## 🔥 Quick Test (CURL)

### 1. Получить токен админа
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@stroymaster.com","password":"admin123"}' | jq -r '.token')

echo $TOKEN
```

### 2. Получить список воркеров
```bash
curl -s http://localhost:8080/api/workers | jq -r '.[0].id'
```

### 3. Продвинуть первого воркера
```bash
WORKER_ID=$(curl -s http://localhost:8080/api/workers | jq -r '.[0].id')

curl -X POST "http://localhost:8080/api/admin/workers/$WORKER_ID/promote" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "promotionType": "premium",
    "durationDays": 30
  }'
```

### 4. Проверить ТОП воркеров
```bash
curl -s http://localhost:8080/api/promotions/top-workers | jq
```

---

## ⚙️ Что нужно добавить на фронтенде

### В админ-панели:
1. **Форма продвижения воркера** на странице воркера
2. **Кнопка "Promote"** с выбором типа и срока
3. **Отображение текущего статуса продвижения**
4. **Кнопка "Cancel Promotion"** для активных продвижений

### На публичной части:
1. **Бейджи** для продвинутых воркеров (PREMIUM/TOP/FEATURED)
2. **Сортировка** - продвинутые воркеры в топе списка
3. **Секция "Top Workers"** на главной странице

### Пример компонента (React):
```jsx
// Admin Panel - Promote Worker Form
<PromotionForm workerId={worker.id}>
  <select name="promotionType">
    <option value="featured">Featured - 100₸/день</option>
    <option value="top">Top - 250₸/день</option>
    <option value="premium">Premium - 500₸/день</option>
  </select>
  <input type="number" name="durationDays" min="7" max="30" />
  <button>Promote Worker</button>
</PromotionForm>

// Public - Worker Badge
{worker.currentPromotion && (
  <Badge type={worker.currentPromotion}>
    {worker.currentPromotion.toUpperCase()}
  </Badge>
)}
```

---

## ✅ Статус

- ✅ База данных (таблицы promotion_pricing, promotion_history)
- ✅ API endpoints (pricing, promote, cancel, history, top-workers)
- ✅ Backend сервисы и хендлеры
- ⏳ Frontend UI (нужно добавить формы и бейджи)
