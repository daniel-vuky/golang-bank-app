package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/daniel-vuky/golang-bank-app/token"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/metadata"

	mockdb "github.com/daniel-vuky/golang-bank-app/db/mock"
	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/pb"
	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestUpdateUserAPI tests the API_UpdateUser API
func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)

	newName := util.RandomOwner()
	newEmail := util.RandomEmail()

	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: pgtype.Text{
						String: newName,
						Valid:  true,
					},
					Email: pgtype.Text{
						String: newEmail,
						Valid:  true,
					},
				}
				updatedUser := db.User{
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					FullName:          newName,
					Email:             newEmail,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
					IsEmailVerified:   user.IsEmailVerified,
				}
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				accessToken, _, err := tokenMaker.CreateToken(user.Username, "", time.Minute)
				require.NoError(t, err)
				bearerToken := fmt.Sprintf("Bearer %s", accessToken)
				md := metadata.MD{
					authorization: []string{bearerToken},
				}
				return metadata.NewIncomingContext(context.Background(), md)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				updatedUser := res.GetUser()
				require.Equal(t, newName, updatedUser.FullName)
				require.Equal(t, newEmail, updatedUser.Email)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)
			server := NewTestServer(t, store, nil)

			ctx := tc.buildContext(t, server.tokenMaker)

			res, err := server.UpdateUser(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
