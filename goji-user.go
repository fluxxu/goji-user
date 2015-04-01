package user

// Schema
// CREATE TABLE `user` (
//   `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
//   `email` varchar(255) NOT NULL,
//   `password` varchar(255) NOT NULL,
//   `display_name` varchar(255) NOT NULL,
//   `created_at` datetime DEFAULT NULL,
//   `updated_at` datetime DEFAULT NULL,
//   PRIMARY KEY (`id`),
//   KEY `user_email_index` (`email`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
// ALTER TABLE `user` ADD UNIQUE INDEX `email_unique` (`email` ASC);

import (
	"github.com/fluxxu/goji-auth"
	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji/web"
)

type Opts struct {
	Dbx     *sqlx.DB
	Mux     *web.Mux
	MuxBase string
}

var opts *Opts
var dbx *sqlx.DB

func Configure(options *Opts) {
	opts = options
	dbx = opts.Dbx

	auth.Skip(opts.MuxBase + "/init")
	opts.Mux.Post(opts.MuxBase+"/init", RouteInit)
}
