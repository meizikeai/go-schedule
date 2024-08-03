package models

import (
	"go-schedule/libs/tool"

	log "github.com/sirupsen/logrus"
)

var tools = tool.NewTools()

// test structure
// CREATE TABLE `test_user_info` (
//   `id` int NOT NULL AUTO_INCREMENT,
//   `email` varchar(50) NOT NULL DEFAULT '' COMMENT '用户帐号',
//   `name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户姓名',
//   `national` varchar(10) NOT NULL DEFAULT '' COMMENT '民族',
//   `gender` varchar(10) NOT NULL DEFAULT '' COMMENT '性别',
//   `idcard` varchar(18) NOT NULL DEFAULT '' COMMENT '身份证号',
//   `phone` varchar(11) NOT NULL DEFAULT '' COMMENT '手机号',
//   `address` varchar(100) NOT NULL DEFAULT '' COMMENT '家庭地址',
//   `postcode` varchar(6) NOT NULL DEFAULT '' COMMENT '邮编号',
//   `datetime` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//   PRIMARY KEY (`id`),
//   UNIQUE KEY `email` (`email`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

type Person struct {
	Id       int    `json:"id" form:"id"`
	Email    string `json:"email" form:"email"`
	Name     string `json:"name" form:"name"`
	National string `json:"national" form:"national"`
	Gender   string `json:"gender" form:"gender"`
	IdCard   string `json:"idcard" form:"idcard"`
	Phone    string `json:"phone" form:"phone"`
	Address  string `json:"address" form:"address"`
	Postcode string `json:"postcode" form:"postcode"`
	Datetime string `json:"datetime" form:"datetime"`
}

func AddPerson(v []string) (id int64, err error) {
	pool := tools.GetMySQLClient("default.master")
	res, err := pool.Exec(`
		INSERT INTO test_user_info(id, email, name, national, gender, idcard, phone, address, postcode)
		VALUES (null, ?, ?, ?, ?, ?, ?, ?, ?)`,
		v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7])

	if err != nil {
		log.Error(err)
	}

	id, err = res.LastInsertId()

	return id, err
}

func GetPerson(email string) (result Person, err error) {
	var person Person

	pool := tools.GetMySQLClient("default.master")
	res := pool.QueryRow("SELECT id, email, name, national, gender, idcard, phone, address, postcode, datetime FROM test_user_info WHERE email=?", email)

	err = res.Scan(
		&person.Id,
		&person.Email,
		&person.Name,
		&person.National,
		&person.Gender,
		&person.IdCard,
		&person.Phone,
		&person.Address,
		&person.Postcode,
		&person.Datetime,
	)

	if err != nil {
		log.Error(err)
	}

	result = person

	return result, err
}

func GetPersons() (result []Person, err error) {
	pool := tools.GetMySQLClient("default.master")
	res, err := pool.Query("SELECT id, email, name, national, gender, idcard, phone, address, postcode, datetime FROM test_user_info")

	if err != nil {
		log.Error(err)
	}

	defer res.Close()

	for res.Next() {
		var person Person
		var datetime string

		res.Scan(
			&person.Id,
			&person.Email,
			&person.Name,
			&person.National,
			&person.Gender,
			&person.IdCard,
			&person.Phone,
			&person.Address,
			&person.Postcode,
			&datetime,
		)

		person.Datetime = datetime
		result = append(result, person)
	}

	return result, err
}

func UpdatePerson(name, phone, email string) (ra int64, err error) {
	pool := tools.GetMySQLClient("default.master")
	row, err := pool.Prepare("UPDATE test_user_info SET name=?, phone=? WHERE id=?")

	defer row.Close()

	res, err := row.Exec(name, phone, email)

	if err != nil {
		log.Error(err)
	}

	ra, err = res.RowsAffected()

	return ra, err
}
