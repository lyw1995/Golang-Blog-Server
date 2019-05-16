package api

// swagger:operation GET /v1/users user getDefaultUser
// 获取默认用户信息
// ---

// swagger:operation GET /v1/users/:uid user getUserById
// 获取指定用户id用户信息
// ---


// swagger:operation POST /v1/upload upload 上传头像/博客文章图片
// 上传头像/博客文章图片
// ---


// swagger:operation GET /v1/users/:uid/links link 获取指定用户所有友链
// 获取指定用户所有友链
// ---

// swagger:operation GET /v1/users/:uid/links/:lid link 获取id友链信息
// 获取id友链信息
// ---



// swagger:operation GET /v1/users/:uid/categorys category 获取指定用户所有分类
// 获取指定用户所有分类
// ---

// swagger:operation GET /v1/users/:uid/categorys/:cid/articles category 获取用户指定分类下所有文章
// 获取用户指定分类下所有文章
// ---





// swagger:operation GET /v1/users/:uid/articles article 获取用户所有文章
// 获取用户所有文章
// ---

// swagger:operation GET /v1/users/:uid/other article 获取热门最新文章
// 获取热门最新文章
// ---

// swagger:operation GET /v1/users/:uid/articles/:aid article 获取某篇文章
// 获取某篇文章
// ---