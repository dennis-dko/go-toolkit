package acl

import (
	"context"
	"testing"

	"github.com/dennis-dko/go-toolkit/testhandler"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type AclTestSuite struct {
	suite.Suite
	ctx      context.Context
	instance *echo.Echo
	config   Config
	userID   string
	roles    []string
}

func (a *AclTestSuite) SetupTest() {
	// Setup
	a.instance = echo.New()
	a.config = Config{
		Enabled:     true,
		Username:    "test",
		Password:    "test",
		AuthModel:   "./testdata/auth.conf",
		PolicyModel: "./testdata/policy.csv",
	}
	a.userID = "test"
	a.roles = []string{
		"admin",
		"user",
		"guest",
	}
}

func (a *AclTestSuite) SetupSubTest() {
	// Sub setup
	a.ctx = testhandler.Ctx(false, false)
}

func TestAclTestSuite(t *testing.T) {
	suite.Run(t, new(AclTestSuite))
}

func (a *AclTestSuite) TestProvide() {

	a.Run("happy path - provide acl enforcer", func() {
		// Run
		err := a.config.Provide()

		// Assert
		a.NoError(err)
	})
	a.Run("should return an error while providing acl enforcer", func() {
		// Init
		a.config.AuthModel = ""

		// Run
		err := a.config.Provide()

		// Assert
		a.Error(err)
	})
}

func (a *AclTestSuite) TestAddUser() {

	a.Run("happy path - add acl user", func() {
		// Run
		err := a.config.Provide()
		addErr := AddUser(a.ctx, a.userID, a.roles)

		// Assert
		a.NoError(err)
		a.NoError(addErr)
	})
}

func (a *AclTestSuite) TestDeleteUser() {

	a.Run("happy path - delete acl user", func() {
		// Run
		err := a.config.Provide()
		addErr := AddUser(a.ctx, a.userID, a.roles)
		deleteErr := DeleteUser(a.ctx, a.userID)

		// Assert
		a.NoError(err)
		a.NoError(addErr)
		a.NoError(deleteErr)
	})
}

func (a *AclTestSuite) TestUserPermissions() {

	a.Run("happy path - get acl user permissions", func() {
		// Run
		err := a.config.Provide()
		addErr := AddUser(a.ctx, a.userID, a.roles)
		perms, permsErr := GetPermissionsForUser(a.ctx, a.userID)

		// Assert
		a.NoError(err)
		a.NoError(addErr)
		a.NoError(permsErr)
		a.NotEmpty(perms)
	})
}

func (a *AclTestSuite) TestAuthRoutes() {

	a.Run("happy path - get authorized acl routes", func() {
		// Run
		err := a.config.Provide()
		authRoutes, authErr := GetAuthorizedRoutes()

		// Assert
		a.NoError(err)
		a.NoError(authErr)
		a.NotEmpty(authRoutes)
	})
}

func (a *AclTestSuite) TestUseEnforcer() {

	a.Run("happy path - use acl enforcer", func() {
		// Run
		err := a.config.Provide()
		UseAuthEnforcer(a.ctx, a.instance)

		// Assert
		a.NoError(err)
	})
}
