package mysql

// database info
const (
	DatabaseName   string = "Cloud"
	createDatabase string = "create database if not exists " + DatabaseName
)

// user table info
const (
	UserTableName string = "userinfo"
	userTable     string = DatabaseName + "." + UserTableName

	insertIntoUserTable string = "insert into " + userTable + " (uname, password, email) values (?, ?, ?) "
	deleteFromUserTable string = "delete from " + userTable + " "
	selectFromUserTable string = "select %s from " + userTable + " "
	updateUserTable     string = "update " + userTable + " set "
	createUserTable     string = `
	create table if not exists ` + userTable + `(
		uid      int unsigned  auto_increment,
		uname    varchar(128)  not null,
		password varchar(128)  not null,
		email    varchar(128)  not null,
		primary key(uid)
	) engine=InnoDB default charset=utf8`
)

// file table info
const (
	FileTableName string = "fileinfo"
	fileTable     string = DatabaseName + "." + FileTableName

	insertIntoFileTable string = "insert into " + fileTable + " (uid, filename, filepath, md5value) values (?, ?, ?, ?) "
	deleteFromFileTable string = "delete from " + fileTable + " "
	selectFromFileTable string = "select %s from " + fileTable + " "
	updateFileTable     string = "update " + fileTable + " set "
	createFileTable     string = `
	create table if not exists ` + fileTable + `(
		fid      int unsigned  auto_increment,
		uid      int unsigned  not null,
		filename varchar(128)  not null,
		filepath varchar(128)  not null,
		md5value varchar(32)   not null,
		primary key(fid),
		foreign key(uid) references userinfo(uid)
	) engine=InnoDB default charset=utf8`
)
