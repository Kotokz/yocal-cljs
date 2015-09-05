package forms

type RegisterForm struct {
	LoginName string `form:"username" binding:"Required;AlphaDashDot;MaxSize(35)"`
	FullName  string `form:"fullname" binding:"Required;MaxSize(35)"`
	Staffid   int    `form:"staffid" binding:"Required;MaxSize(10)"`
	Email     string `form:"email" binding:"Required;Email;MaxSize(50)"`
	Password  string `form:"password" binding:"Required;MinSize(6);MaxSize(255)"`
	Retype    string `form:"retype" binding:"Required;MinSize(6);MaxSize(255)"`
}

type LoginForm struct {
	LoginName string `form:"username" binding:"Required;MaxSize(35)"`
	Password  string `form:"password" binding:"Required;MinSize(6);MaxSize(255)"`
	//Remember  bool   `form:"remember"`
}
