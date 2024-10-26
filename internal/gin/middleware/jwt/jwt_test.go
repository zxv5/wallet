package jwt

// // Mocking the JWT library
// type MockJWT struct {
// 	mock.Mock
// }

// func (m *MockJWT) Sign() (string, error) {
// 	args := m.Called()
// 	return args.String(0), args.Error(1)
// }

// func (m *MockJWT) Parse(token string) (*jwtv5.Token, error) {
// 	args := m.Called(token)
// 	return args.Get(0).(*jwtv5.Token), args.Error(1)
// }

// // MockUserInfo for user information
// var mockUserInfo = types.UserInfo{
// 	ID: 1,
// }

// // setupRouter initializes the Gin router with JWT middleware
// func setupRouter() *gin.Engine {
// 	gin.SetMode(gin.TestMode)
// 	r := gin.New()
// 	r.Use(Jwt())
// 	return r
// }

// func TestJwtMiddlewareSuccess(t *testing.T) {
// 	config.Config.Jwt.Secret = "testsecret"
// 	config.Config.Jwt.Exp = 1 // 1 hour expiration

// 	// Create a valid token using the Sign function
// 	token, err := Sign(&mockUserInfo)
// 	assert.NoError(t, err)

// 	router := setupRouter()
// 	router.GET("/protected", func(c *gin.Context) {
// 		userInfo, exists := c.Get(UserInfoKey)
// 		assert.True(t, exists)
// 		assert.Equal(t, mockUserInfo, userInfo)
// 		c.JSON(http.StatusOK, gin.H{"message": "success"})
// 	})

// 	req, _ := http.NewRequest("GET", "/protected", nil)
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	resp := httptest.NewRecorder()
// 	router.ServeHTTP(resp, req)

// 	assert.Equal(t, http.StatusOK, resp.Code)
// 	assert.JSONEq(t, `{"message": "success"}`, resp.Body.String())
// }

// func TestJwtMiddlewareInvalidToken(t *testing.T) {
// 	router := setupRouter()
// 	router.GET("/protected", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"message": "success"})
// 	})

// 	req, _ := http.NewRequest("GET", "/protected", nil)
// 	req.Header.Set("Authorization", "Bearer invalidtoken")
// 	resp := httptest.NewRecorder()
// 	router.ServeHTTP(resp, req)

// 	assert.Equal(t, http.StatusUnauthorized, resp.Code)
// }

// func TestJwtMiddlewareMissingToken(t *testing.T) {
// 	router := setupRouter()
// 	router.GET("/protected", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"message": "success"})
// 	})

// 	req, _ := http.NewRequest("GET", "/protected", nil)
// 	resp := httptest.NewRecorder()
// 	router.ServeHTTP(resp, req)

// 	assert.Equal(t, http.StatusUnauthorized, resp.Code)
// }
