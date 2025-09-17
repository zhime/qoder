package api

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

// RefreshRequest 刷新令牌请求
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// User 用户信息（用于响应）
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   int    `json:"status"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=user admin"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

// CreateServerRequest 创建服务器请求
type CreateServerRequest struct {
	Name        string `json:"name" binding:"required"`
	Host        string `json:"host" binding:"required"`
	Port        int    `json:"port" binding:"required,min=1,max=65535"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password"`
	PrivateKey  string `json:"private_key"`
	Environment string `json:"environment" binding:"required,oneof=dev test prod"`
	Description string `json:"description"`
}

// UpdateServerRequest 更新服务器请求
type UpdateServerRequest struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        *int   `json:"port" binding:"omitempty,min=1,max=65535"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PrivateKey  string `json:"private_key"`
	Environment string `json:"environment" binding:"omitempty,oneof=dev test prod"`
	Description string `json:"description"`
	Status      *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

// CreateDeploymentRequest 创建部署请求
type CreateDeploymentRequest struct {
	Name       string `json:"name" binding:"required"`
	ServerID   uint   `json:"server_id" binding:"required"`
	Repository string `json:"repository" binding:"required"`
	Branch     string `json:"branch"`
	Path       string `json:"path" binding:"required"`
	Script     string `json:"script" binding:"required"`
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Name     string `json:"name" binding:"required"`
	Command  string `json:"command" binding:"required"`
	CronExpr string `json:"cron_expr" binding:"required"`
	ServerID uint   `json:"server_id" binding:"required"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cron_expr"`
	ServerID *uint  `json:"server_id"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

// PageRequest 分页请求
type PageRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=10" binding:"min=1,max=100"`
}

// PageResponse 分页响应
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}