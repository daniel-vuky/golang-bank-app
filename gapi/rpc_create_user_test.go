package gapi

import (
	"context"
	"fmt"
	mockdb "github.com/daniel-vuky/golang-bank-app/db/mock"
	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/pb"
	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/daniel-vuky/golang-bank-app/worker"
	mockWorker "github.com/daniel-vuky/golang-bank-app/worker/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (e eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	if !reflect.DeepEqual(e.arg.CreateUserParams, arg.CreateUserParams) {
		return false
	}
	err = arg.AfterCreate(e.user)
	return err == nil
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

// randomUser get a random user
func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hashedPassword, hashedPasswordErr := util.HashPassword(password)
	require.NoError(t, hashedPasswordErr)

	user := db.User{
		Username:       util.RandomString(8),
		FullName:       util.RandomString(8),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
	}

	return user, password
}

// TestAPI_CreateUser tests the API_CreateUser API
func TestAPI_CreateUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockWorker.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockWorker.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username:       user.Username,
						FullName:       user.FullName,
						HashedPassword: user.HashedPassword,
						Email:          user.Email,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{
						User: user,
					}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUser := res.GetUser()
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.FullName, createdUser.FullName)
				require.Equal(t, user.Email, createdUser.Email)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			taskController := gomock.NewController(t)
			defer taskController.Finish()
			taskDistributor := mockWorker.NewMockTaskDistributor(taskController)

			tc.buildStubs(store, taskDistributor)
			server := NewTestServer(t, store, taskDistributor)

			res, err := server.CreateUser(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
