package api

// swagger:operation GET /admin/users user 获取所有用户列表
// 获取所有用户列表
// ---

// swagger:operation GET /admin/users/:uid user 获取某个用户信息
// 获取某个用户信息
// ---

// swagger:operation POST /admin/users user 创建用户
// 创建用户
// ---

// swagger:operation DELETE /admin/users/:uid user 删除用户
// 删除用户
// ---

// swagger:operation PUT /admin/users/:uid user 修改用户信息
// 修改用户信息
// ---

// swagger:operation POST /v1/admin/sessions session 管理员登录
// 管理员登录
// ---

// swagger:operation POST /v1/admin/sessions/upload session 管理员登录上传图片信息
// 管理员登录上传图片信息
// ---

// swagger:operation DELETE /v1/admin/sessions/:uid session 退出登录
// 退出登录
// ---

// swagger:operation GET /v1/admin/extends extend 获取网站统计信息
// 获取网站统计信息
// ---

// swagger:operation POST /v1/admin/extends extend 传入文章url采集csdn文章
// 传入文章url采集csdn文章
// ---

// swagger:operation GET /admin/users/:uid/links link 获取用户所有友链
// 获取用户所有友链
// ---

// swagger:operation GET /admin/users/:uid/links/:lid link 获取指定友链信息
// 获取指定友链信息
// ---

// swagger:operation POST /admin/users/:uid/links link 添加友链
// 添加友链
// ---

// swagger:operation DELETE /admin/users/:uid/links/:lid link 删除友链
// 删除友链
// ---

// swagger:operation PUT /admin/users/:uid/links/:lid link 修改友链信息
// 修改友链信息
// ---


// swagger:operation GET /admin/users/:uid/categorys category 获取用户所有分类
// 获取用户所有分类
// ---

// swagger:operation GET /admin/users/:uid/categorys/:cid category 获取指定分类信息
// 获取指定分类信息
// ---

// swagger:operation POST /admin/users/:uid/categorys category 添加分类
// 添加分类
// ---

// swagger:operation DELETE /admin/users/:uid/categorys/:cid category 删除分类
// 删除分类
// ---

// swagger:operation PUT /admin/users/:uid/categorys/:cid category 修改分类信息
// 修改分类信息
// ---

// swagger:operation POST /admin/users/:uid/categorys/:cid/articles article 在某分类下创建文章
// 在某分类下创建文章
// ---

// swagger:operation PUT /admin/users/:uid/categorys/:cid/articles/:aid article 修改文章
// 修改文章
// ---

// swagger:operation DELETE /admin/users/:uid/categorys/:cid/articles article 删除分类下所有文章
// 删除分类下所有文章
// ---

// swagger:operation DELETE /admin/users/:uid/categorys/:cid/articles/:aid article 删除某篇文章
// 删除某篇文章
// ---