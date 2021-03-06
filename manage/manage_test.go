package manage_test

import (
	"testing"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/store"

	. "github.com/smartystreets/goconvey/convey"
)

func TestManager(t *testing.T) {
	Convey("Manager test", t, func() {
		manager := manage.NewDefaultManager()
		manager.MapClientStorage(store.NewTestClientStore())
		manager.MustTokenStorage(store.NewMemoryTokenStore())

		Convey("CheckInterface test", func() {
			err := manager.CheckInterface()
			So(err, ShouldBeNil)
		})

		Convey("GetClient test", func() {
			cli, err := manager.GetClient("1")
			So(err, ShouldBeNil)
			So(cli.GetSecret(), ShouldEqual, "11")
		})

		Convey("Token test", func() {
			testManager(manager)
		})
	})
}

func testManager(manager oauth2.Manager) {
	reqParams := &oauth2.TokenGenerateRequest{
		ClientID:    "1",
		UserID:      "123456",
		RedirectURI: "http://localhost/oauth2",
		Scope:       "all",
	}
	cti, err := manager.GenerateAuthToken(oauth2.Code, reqParams)
	So(err, ShouldBeNil)

	code := cti.GetCode()
	So(code, ShouldNotBeEmpty)

	atParams := &oauth2.TokenGenerateRequest{
		ClientID:     reqParams.ClientID,
		ClientSecret: "11",
		RedirectURI:  reqParams.RedirectURI,
		Code:         code,
	}
	ati, err := manager.GenerateAccessToken(oauth2.AuthorizationCode, atParams)
	So(err, ShouldBeNil)

	accessToken, refreshToken := ati.GetAccess(), ati.GetRefresh()
	So(accessToken, ShouldNotBeEmpty)
	So(refreshToken, ShouldNotBeEmpty)

	ainfo, err := manager.LoadAccessToken(accessToken)
	So(err, ShouldBeNil)
	So(ainfo.GetClientID(), ShouldEqual, atParams.ClientID)

	rinfo, err := manager.LoadRefreshToken(refreshToken)
	So(err, ShouldBeNil)
	So(rinfo.GetClientID(), ShouldEqual, atParams.ClientID)

	atParams.Refresh = refreshToken
	atParams.Scope = "owner"
	rti, err := manager.RefreshAccessToken(atParams)
	So(err, ShouldBeNil)

	refreshAT := rti.GetAccess()
	So(refreshAT, ShouldNotBeEmpty)

	_, err = manager.LoadAccessToken(accessToken)
	So(err, ShouldNotBeNil)

	refreshAInfo, err := manager.LoadAccessToken(refreshAT)
	So(err, ShouldBeNil)
	So(refreshAInfo.GetScope(), ShouldEqual, "owner")

	err = manager.RemoveAccessToken(refreshAT)
	So(err, ShouldBeNil)

	_, err = manager.LoadAccessToken(refreshAT)
	So(err, ShouldNotBeNil)

	err = manager.RemoveRefreshToken(refreshToken)
	So(err, ShouldBeNil)

	_, err = manager.LoadRefreshToken(refreshToken)
	So(err, ShouldNotBeNil)
}
