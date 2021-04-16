package models

import (
	"go-schedule/libs/tool"

	log "github.com/sirupsen/logrus"
)

type LoverGift struct {
	Uid   int    `json:"uid" form:"uid"`
	Beans string `json:"beans" form:"beans"`
}

// 数据结构
// CREATE TABLE `lover_gift` (
//   `id` int NOT NULL AUTO_INCREMENT,
//   `uid` bigint NOT NULL,
//   `target_uid` bigint NOT NULL,
//   `live_id` int NOT NULL,
//   `goods_id` bigint NOT NULL,
//   `pay_time` int NOT NULL,
//   `beans` float NOT NULL,
//   `discount_id` bigint NOT NULL,
//   `status` int NOT NULL,
//   `platform` varchar(10) NOT NULL,
//   `consumes_id` varchar(32) NOT NULL,
//   `count` int NOT NULL DEFAULT '1',
//   `ops` int NOT NULL DEFAULT '0',
//   `day` int NOT NULL,
//   `datetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
//   PRIMARY KEY (`id`),
//   UNIQUE KEY `id` (`id`),
//   KEY `uid` (`uid`),
//   KEY `target_uid` (`target_uid`),
//   KEY `day_goods` (`day`,`goods_id`,`platform`),
//   KEY `day_platform` (`day`,`platform`),
//   KEY `target_goods` (`target_uid`,`platform`,`goods_id`),
//   KEY `uid_goods` (`uid`,`platform`,`goods_id`),
//   KEY `target_platform` (`target_uid`,`platform`),
//   KEY `uid_platform` (`uid`,`platform`)
// ) ENGINE=InnoDB AUTO_INCREMENT=1339 DEFAULT CHARSET=utf8;

// 使用的源
// INSERT INTO lover_gift (id, uid, target_uid, live_id, goods_id, pay_time, beans, discount_id, status, platform, consumes_id, count, ops, day) VALUES (null,15410237,15411714,719272,1401,1580880680,10,0,1,"ios","0",1,0,"20200205")

func GetLoverGift(day string) (result []LoverGift, err error) {
	pool := tool.GetMySQLClient("default.slave")

	result = make([]LoverGift, 0)

	sql := "SELECT target_uid as uid,beans FROM lover_gift	WHERE status = 1 AND day = " + day

	rows, err := pool.Query(sql)

	if err != nil {
		log.Error(err)

		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var person LoverGift

		rows.Scan(
			&person.Uid,
			&person.Beans,
		)

		result = append(result, person)
	}

	return result, nil
}

// 数据结构
// CREATE TABLE `lover_daily_beans` (
//   `uid` int NOT NULL DEFAULT '0',
//   `beans` int DEFAULT '0',
//   `day` int NOT NULL DEFAULT '0',
//   `updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
//   PRIMARY KEY (`uid`,`day`),
//   KEY `day` (`day`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

func UpdateLoverData(t string, v string, u string) (id int64, err error) {
	pool := tool.GetMySQLClient("default.master")

	sql := "INSERT INTO " + t + " (uid,beans,day) VALUES " + v + " ON DUPLICATE KEY UPDATE " + u
	// log.Info(sql)

	res, err := pool.Exec(sql)

	if err != nil {
		log.Error(err)
	}

	id, err = res.RowsAffected()

	return id, err
}
