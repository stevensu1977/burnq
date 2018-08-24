package service

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"

	"github.com/stevensu1977/burnq/model"
	"github.com/stevensu1977/toolbox/crypto"

	log "github.com/sirupsen/logrus"
)

var db *storm.DB

func init() {
	var err error

	if _, err = os.Stat("data/"); os.IsNotExist(err) {
		os.Mkdir("./data", os.ModePerm)
	}

	db, err = storm.Open("burnq.db")
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Init db success")
	if !HaveAdmin() {
		log.Info("Please setup a admin for burnq")
		SetupAdmin()
	}
}

//Close close storm DB instance
func Close() {
	db.Close()
}

//SetupAdmin init admin/password
func SetupAdmin() {

emailLabel:
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Email: ")
	email, _ := reader.ReadString('\n')
	if email == "\n" || strings.Index(email, "@") == -1 {
		goto emailLabel
	}
passLabel:
	fmt.Println("Enter Password (length 8-16): ")
	password, _ := reader.ReadString('\n')
	if password == "\n" || len(strings.TrimSpace(password)) < 8 || len(strings.TrimSpace(password)) > 16 {
		goto passLabel
	}
	fmt.Println("Enter Confrim Password (length 8-16): ")
	password1, _ := reader.ReadString('\n')
	if password != password1 {
		goto passLabel
	}
	admin := &model.Admin{Email: strings.TrimSpace(email), Password: strings.TrimSpace(password)}

	err := SaveAdmin(admin)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("init admin successful")
	return

}

//SaveAdmin insert admin to db
func SaveAdmin(admin *model.Admin) error {

	password, err := crypto.AESEncrypt(admin.Password)
	if err != nil {
		return err
	}
	admin.Password = password
	if admin.Password == "" {
		return fmt.Errorf("Encrypt failure")
	}
	return db.From("admin").Save(admin)

}

//HaveAdmin check admin account isexits
func HaveAdmin() bool {
	var admins []model.Admin
	db.From("admin").All(&admins)
	if len(admins) > 0 {
		for idx, _ := range admins {
			//fmt.Printf("|%s|%s", admins[idx].Email, admins[idx].Password)
			log.Debugf("|%s|%s", admins[idx].Email, admins[idx].Password)
		}
		return true
	}
	return false
}

//RemoveAdmin drop admin bucket
func RemoveAdmin() {
	err := db.Drop("admin")
	fmt.Println(err)
}

//FetchAdmin load admin by email/passwd
func FetchAdmin(email, password string) (*model.Admin, error) {

	var admin model.Admin
	err := db.From("admin").One("Email", email, &admin)
	if err != nil {
		log.Println("xxx")
		return nil, err
	}

	_password, err := crypto.AESEncrypt(password)
	if err != nil {
		return nil, err
	}

	if admin.Password != _password {
		return nil, fmt.Errorf("Password not match")
	}

	return &admin, nil
}

//UpdateAdminPasswd update admin
func UpdateAdminPasswd(admin *model.Admin) error {
	password, err := crypto.AESEncrypt(admin.Password)
	if err != nil {
		return err
	}
	admin.Password = password
	return db.From("admin").Update(admin)
}

//SaveAccount save cloud account to db
func SaveAccount(account *model.CloudAccount) error {
	password, err := crypto.AESEncrypt(account.Password)
	if err != nil {
		return err
	}
	account.Password = password
	return db.From("user").Save(account)
}

//FetchAccount load cloud account by tenant
func FetchAccount(tenant string, decrypt bool) (*model.CloudAccount, error) {
	var account model.CloudAccount
	err := db.From("user").One("Tenant", tenant, &account)
	if err != nil {
		return nil, err
	}

	if decrypt {
		password, err := crypto.AESDecrypt(account.Password)
		if err != nil {
			return nil, err
		}
		account.Password = password

	}
	return &account, nil
}

//FetchAllAccount load all account
func FetchAllAccount() ([]model.CloudAccount, error) {

	var accounts []model.CloudAccount
	err := db.From("user").All(&accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

//RemoveAccount delete cloud account from db
func RemoveAccount(id int) error {
	query := db.From("user").Select(q.Eq("ID", id))
	return query.Delete(new(model.CloudAccount))
}

//FetchAllTenant just return all tenant
func FetchAllTenant() ([]string, error) {

	var accounts []model.CloudAccount
	var tenant []string
	err := db.From("user").All(&accounts)
	if err != nil {
		return nil, err
	}

	for i, _ := range accounts {
		tenant = append(tenant, accounts[i].Tenant)
	}
	return tenant, nil
}

//UpdateAccount update cloud account info
func UpdateAccount(account *model.CloudAccount) error {
	return db.From("user").Update(account)
}
