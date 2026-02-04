package user_test

import (
	"testing"

	"github.com/ap1-final-mini-moodle/internal/domain/user"
	userUC "github.com/ap1-final-mini-moodle/internal/usecase/user"
)

type mockUserRepo struct {
	users map[string]*user.User
}

func (m *mockUserRepo) Create(u *user.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) GetByID(id string) (*user.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockUserRepo) GetByEmail(email string) (*user.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) Update(u *user.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) List(skip, take int) ([]*user.User, error) {
	return nil, nil
}

func (m *mockUserRepo) Delete(id string) error {
	delete(m.users, id)
	return nil
}

func TestRegisterUser(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*user.User)}
	service := userUC.NewService(repo)

	input := &user.CreateUserInput{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      user.RoleStudent,
	}

	u, err := service.Register(input)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	if u.Email != input.Email {
		t.Errorf("expected email %s, got %s", input.Email, u.Email)
	}
	if u.FirstName != input.FirstName {
		t.Errorf("expected firstName %s, got %s", input.FirstName, u.FirstName)
	}
	if !u.Active {
		t.Errorf("expected user to be active")
	}
}

func TestUpdateUser(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*user.User)}
	service := userUC.NewService(repo)

	// Create user first
	input := &user.CreateUserInput{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      user.RoleStudent,
	}
	u, _ := service.Register(input)

	// Update user
	newFirstName := "Jane"
	updateInput := &user.UpdateUserInput{
		FirstName: &newFirstName,
	}

	updated, err := service.Update(u.ID, updateInput)
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	if updated.FirstName != newFirstName {
		t.Errorf("expected firstName %s, got %s", newFirstName, updated.FirstName)
	}
}

func TestGetUserByID(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*user.User)}
	service := userUC.NewService(repo)

	input := &user.CreateUserInput{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      user.RoleStudent,
	}

	created, _ := service.Register(input)
	retrieved, err := service.GetByID(created.ID)

	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	if retrieved.Email != created.Email {
		t.Errorf("expected email %s, got %s", created.Email, retrieved.Email)
	}
}
