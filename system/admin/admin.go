// Package admin desrcibes the admin view containing references to
// various managers and editors
package admin

import (
	"bytes"
	"encoding/json"
	"github.com/monstrum/ponzu-cms/system/twig"
	"html/template"
	"net/http"

	"github.com/monstrum/ponzu-cms/system/admin/user"
	"github.com/monstrum/ponzu-cms/system/api/analytics"
	"github.com/monstrum/ponzu-cms/system/db"
	"github.com/monstrum/ponzu-cms/system/item"
)

var startAdminHTML = `<!doctype html>
<html lang="en">
    <head>
        <title>{{ .Logo }}</title>
        <script type="text/javascript" src="/admin/static/common/js/jquery-2.1.4.min.js"></script>
        <script type="text/javascript" src="/admin/static/common/js/util.js"></script>
        <script type="text/javascript" src="/admin/static/dashboard/js/materialize.min.js"></script>
        <script type="text/javascript" src="/admin/static/dashboard/js/chart.bundle.min.js"></script>
        <script type="text/javascript" src="/admin/static/editor/js/materialNote.js"></script> 
        <script type="text/javascript" src="/admin/static/editor/js/ckMaterializeOverrides.js"></script>
                  
        <link rel="stylesheet" href="/admin/static/dashboard/css/material-icons.css" />     
        <link rel="stylesheet" href="/admin/static/dashboard/css/materialize.min.css" />
        <link rel="stylesheet" href="/admin/static/editor/css/materialNote.css" />
        <link rel="stylesheet" href="/admin/static/dashboard/css/admin.css" />    

        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
    </head>
    <body class="grey lighten-4">
       <div class="navbar-fixed">
            <nav class="grey darken-2">
            <div class="nav-wrapper">
                <a class="brand-logo" href="/admin">{{ .Logo }}</a>

                <ul class="right">
                    <li><a href="/admin/logout">Logout</a></li>
                </ul>
            </div>
            </nav>
        </div>

        <div class="admin-ui row">`

var endAdminHTML = `
        </div>
        <footer class="row">
            <div class="col s12">
                <p class="center-align">Powered by &copy; <a target="_blank" href="https://ponzu-cms.org">Ponzu</a> &nbsp;&vert;&nbsp; open-sourced by <a target="_blank" href="https://www.bosssauce.it">Boss Sauce Creative</a></p>
            </div>     
        </footer>
    </body>
</html>`

type admin struct {
	Logo    string
	Types   map[string]func() interface{}
	Subview template.HTML
}

// GetDefaultData ...
func GetDefaultData() map[string]interface{} {
	cfg, err := db.Config("name")
	if err != nil {
		return map[string]interface{}{}
	}

	if cfg == nil {
		cfg = []byte("")
	}

	return map[string]interface{}{
		"Logo": string(cfg),
	}
}

// Admin ...
func Admin(view []byte) (_ []byte, err error) {
	cfg, err := db.Config("name")
	if err != nil {
		return
	}

	if cfg == nil {
		cfg = []byte("")
	}

	return twig.Twig.Render("admin.html.twig", map[string]interface{}{
		"Logo":    string(cfg),
		"Types":   item.Types,
		"Subview": string(view),
	})
}

// Init ...
func Init() ([]byte, error) {
	return twig.Twig.Render("init.html.twig", GetDefaultData())
}

// Login ...
func Login() ([]byte, error) {
	return twig.Twig.Render("login.html.twig", GetDefaultData())
}

// ForgotPassword ...
func ForgotPassword() ([]byte, error) {
	return twig.Twig.Render("forgot-password.html.twig", GetDefaultData())
}

// RecoveryKey ...
func RecoveryKey() ([]byte, error) {
	return twig.Twig.Render("recovery-password.html.twig", GetDefaultData())
}

// UsersList ...
func UsersList(req *http.Request) ([]byte, error) {
	html := `
    <div class="card user-management">
        <div class="card-title">Edit your account:</div>    
        <form class="row" enctype="multipart/form-data" action="/admin/configure/users/edit" method="post">
            <div class="col s9">
                <label class="active">Email Address</label>
                <input type="email" name="email" value="{{ .User.Email }}"/>
            </div>

            <div class="col s9">
                <div>To approve changes, enter your password:</div>
                
                <label class="active">Current Password</label>
                <input type="password" name="password"/>
            </div>

            <div class="col s9">
                <label class="active">New Password: (leave blank if no password change needed)</label>
                <input name="new_password" type="password"/>
            </div>

            <div class="col s9">                        
                <button class="btn waves-effect waves-light green right" type="submit">Save</button>
            </div>
        </form>

        <div class="card-title">Add a new user:</div>        
        <form class="row" enctype="multipart/form-data" action="/admin/configure/users" method="post">
            <div class="col s9">
                <label class="active">Email Address</label>
                <input type="email" name="email" value=""/>
            </div>

            <div class="col s9">
                <label class="active">Password</label>
                <input type="password" name="password"/>
            </div>

            <div class="col s9">            
                <button class="btn waves-effect waves-light green right" type="submit">Add User</button>
            </div>   
        </form>        

        <div class="card-title">Remove Admin Users</div>        
        <ul class="users row">
            {{ range .Users }}
            <li class="col s9">
                {{ .Email }}
                <form enctype="multipart/form-data" class="delete-user __ponzu right" action="/admin/configure/users/delete" method="post">
                    <span>Delete</span>
                    <input type="hidden" name="email" value="{{ .Email }}"/>
                    <input type="hidden" name="id" value="{{ .ID }}"/>
                </form>
            </li>
            {{ end }}
        </ul>
    </div>
    `
	script := `
    <script>
        $(function() {
            var del = $('.delete-user.__ponzu span');
            del.on('click', function(e) {
                if (confirm("[Ponzu] Please confirm:\n\nAre you sure you want to delete this user?\nThis cannot be undone.")) {
                    $(e.target).parent().submit();
                }
            });
        });
    </script>
    `
	// get current user out to pass as data to execute template
	j, err := db.CurrentUser(req)
	if err != nil {
		return nil, err
	}

	var usr user.User
	err = json.Unmarshal(j, &usr)
	if err != nil {
		return nil, err
	}

	// get all users to list
	jj, err := db.UserAll()
	if err != nil {
		return nil, err
	}

	var usrs []user.User
	for i := range jj {
		var u user.User
		err = json.Unmarshal(jj[i], &u)
		if err != nil {
			return nil, err
		}
		if u.Email != usr.Email {
			usrs = append(usrs, u)
		}
	}

	// make buffer to execute html into then pass buffer's bytes to Admin
	buf := &bytes.Buffer{}
	tmpl := template.Must(template.New("users").Parse(html + script))
	data := map[string]interface{}{
		"User":  usr,
		"Users": usrs,
	}

	err = tmpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return Admin(buf.Bytes())
}

// Dashboard returns the admin view with analytics dashboard
func Dashboard() ([]byte, error) {
	data, err := analytics.ChartData()
	if err != nil {
		return nil, err
	}
	return twig.Twig.Render("analytics.html.twig", data)
}

func ErrorPage(err string) ([]byte, error) {
	return twig.Twig.Render("error/"+err+".html.twig", GetDefaultData())
}

// Error400 creates a subview for a 400 error page
func Error400() ([]byte, error) {
	return ErrorPage("400")
}

// Error404 creates a subview for a 404 error page
func Error404() ([]byte, error) {
	return ErrorPage("404")
}

// Error405 creates a subview for a 405 error page
func Error405() ([]byte, error) {
	return ErrorPage("405")
}

// Error500 creates a subview for a 500 error page
func Error500() ([]byte, error) {
	return ErrorPage("500")
}

// ErrorMessage is a generic error message container, similar to Error500() and
// others in this package, ecxept it expects the caller to provide a title and
// message to describe to a view why the error is being shown
func ErrorMessage(title, message string) ([]byte, error) {
	data := GetDefaultData()
	data["title"] = title
	data["message"] = message
	return twig.Twig.Render("error/default.html.twig", data)
}
