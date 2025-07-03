package patient

import (
	"errors"
	"go-gin-auth/utils"
	"strings"
	"time"
)

var (
	ErrNotFound       = errors.New("pasien tidak ditemukan")
	ErrInvalidInput   = errors.New("input tidak valid atau tidak lengkap")
	ErrIdentityExists = errors.New("nomor identitas sudah digunakan oleh pasien lain yang aktif")
)

type Service interface {
	CreatePatient(input *Patient) (*Patient, error)
	GetAllPatients(searchQuery string) ([]Patient, error)
	GetPatientByID(id uint) (*Patient, error)
	UpdatePatient(id uint, input *Patient) (*Patient, error)
	DeletePatient(id uint) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) encryptPatientData(p *Patient, dateOfBirth string) error {
	var err error
	p.FullName, err = utils.Encrypt(p.FullName)
	if err != nil {
		return err
	}
	if dateOfBirth != "" {
		p.DateOfBirth, err = utils.Encrypt(dateOfBirth)
		if err != nil {
			return err
		}
	}
	p.Address, err = utils.Encrypt(p.Address)
	if err != nil {
		return err
	}
	p.PhoneNumber, err = utils.Encrypt(p.PhoneNumber)
	if err != nil {
		return err
	}
	p.IdentityNumber, err = utils.Encrypt(p.IdentityNumber)
	if err != nil {
		return err
	}
	if p.GuarantorName != "" {
		p.GuarantorName, err = utils.Encrypt(p.GuarantorName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) decryptPatientData(p *Patient) error {
	var err error
	p.FullName, err = utils.Decrypt(p.FullName)
	if err != nil {
		return err
	}
	if p.DateOfBirth != "" {
		decryptedDOB, err := utils.Decrypt(p.DateOfBirth)
		if err != nil {
			return err
		}
		p.DecryptedDateOfBirth = decryptedDOB
		dob, _ := time.Parse("2006-01-02", decryptedDOB)
		if !dob.IsZero() {
			now := time.Now()
			age := now.Year() - dob.Year()
			if now.YearDay() < dob.YearDay() {
				age--
			}
			p.Age = age
		}
	}
	p.Address, err = utils.Decrypt(p.Address)
	if err != nil {
		return err
	}
	p.PhoneNumber, err = utils.Decrypt(p.PhoneNumber)
	if err != nil {
		return err
	}
	p.IdentityNumber, err = utils.Decrypt(p.IdentityNumber)
	if err != nil {
		return err
	}
	if p.GuarantorName != "" {
		p.GuarantorName, err = utils.Decrypt(p.GuarantorName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) CreatePatient(input *Patient) (*Patient, error) {
	input.FullName = strings.TrimSpace(input.FullName)
	if input.FullName == "" || input.Gender == "" || input.PlaceOfBirth == "" || input.DecryptedDateOfBirth == "" {
		return nil, ErrInvalidInput
	}

	if err := s.encryptPatientData(input, input.DecryptedDateOfBirth); err != nil {
		return nil, errors.New("gagal mengenkripsi data pasien")
	}

	if input.Status == "" {
		input.Status = "Aktif"
	}

	newPatient, err := s.repository.Create(input)
	if err != nil {
		return nil, err
	}

	if err := s.decryptPatientData(newPatient); err != nil {
		return nil, errors.New("gagal mendekripsi data untuk respons")
	}
	return newPatient, nil
}

func (s *service) UpdatePatient(id uint, input *Patient) (*Patient, error) {
	_, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.encryptPatientData(input, input.DecryptedDateOfBirth); err != nil {
		return nil, errors.New("gagal mengenkripsi data untuk update")
	}

	updatedPatient, err := s.repository.Update(id, input)
	if err != nil {
		return nil, err
	}

	if err := s.decryptPatientData(updatedPatient); err != nil {
		return nil, errors.New("gagal mendekripsi data untuk respons")
	}
	return updatedPatient, nil
}

func (s *service) GetAllPatients(searchQuery string) ([]Patient, error) {
	patients, err := s.repository.GetAll(searchQuery)
	if err != nil {
		return nil, err
	}

	decryptedPatients := make([]Patient, 0)
	for _, p := range patients {
		if err := s.decryptPatientData(&p); err == nil {
			decryptedPatients = append(decryptedPatients, p)
		}
	}

	if searchQuery != "" {
		var filtered []Patient
		searchQuery = strings.ToLower(searchQuery)
		for _, p := range decryptedPatients {
			if strings.Contains(strings.ToLower(p.FullName), searchQuery) || strings.Contains(p.IdentityNumber, searchQuery) {
				filtered = append(filtered, p)
			}
		}
		return filtered, nil
	}
	return decryptedPatients, nil
}

func (s *service) GetPatientByID(id uint) (*Patient, error) {
	patient, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.decryptPatientData(patient); err != nil {
		return nil, errors.New("gagal mendekripsi data pasien")
	}
	return patient, nil
}

func (s *service) DeletePatient(id uint) error {
	if _, err := s.repository.GetByID(id); err != nil {
		return err
	}
	return s.repository.Delete(id)
}
