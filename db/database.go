package db

import (
	"encoding/json"
	"errors"
	"medicontrol/auth"
	"medicontrol/models"
	"os"
	"sync"
	"time"
)

// Este arquivo é um stub para satisfazer as dependências de compilação.
// A lógica real do banco de dados está no pacote 'models' usando SQLite.

type Database struct {
	Users         []auth.User           `json:"users"`
	Medicamentos  []models.Medicamento  `json:"medicamentos"`
	Movimentacoes []models.Movimentacao `json:"movimentacoes"`
	mutex         sync.RWMutex
	usersFile     string
	medsFile      string
}

func (d *Database) CreateMedicamento(med models.Medicamento) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	med.CriadoEm = time.Now()

	d.Medicamentos = append(d.Medicamentos, med)
	return d.save()
}

func (d *Database) UpdateMedicamento(med models.Medicamento) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for i, m := range d.Medicamentos {
		if m.ID == med.ID {
			med.CriadoEm = time.Now()
			d.Medicamentos[i] = med
			return d.save()
		}
	}

	return errors.New("medicamento não encontrado")
}

func (d *Database) RegistrarMovimentacao(mov models.Movimentacao) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	mov.Data = time.Now()

	// Atualiza quantidade do medicamento
	for i, med := range d.Medicamentos {
		if med.ID == mov.MedicamentoID {
			if mov.Tipo == "entrada" {
				d.Medicamentos[i].Quantidade += mov.Quantidade
			} else if mov.Tipo == "saida" {
				if d.Medicamentos[i].Quantidade < mov.Quantidade {
					return errors.New("quantidade insuficiente em estoque")
				}
				d.Medicamentos[i].Quantidade -= mov.Quantidade
			}
			d.Medicamentos[i].CriadoEm = time.Now()
			break
		}
	}

	d.Movimentacoes = append(d.Movimentacoes, mov)
	return d.save()
}

var db *Database

func Init() error {
	db = &Database{
		usersFile: "data/users.json",
	}
	return db.load()
}

func (d *Database) load() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	data, err := os.ReadFile(d.usersFile)
	if os.IsNotExist(err) {
		return d.save()
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &d)
}

func (d *Database) save() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	data, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(d.usersFile, data, 0644)
}

func GetUserByUsername(username string) (auth.User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, user := range db.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return auth.User{}, errors.New("usuário não encontrado")
}

func CreateUser(user auth.User) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Verifica se usuário já existe
	for _, u := range db.Users {
		if u.Username == user.Username {
			return errors.New("usuário já existe")
		}
	}

	// Gera ID para o novo usuário
	user.ID = len(db.Users) + 1
	db.Users = append(db.Users, user)

	return db.save()
}
