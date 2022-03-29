package user_test

import (
	"regexp"
	"testing"

	um "game-api/pkg/domain/model/user"
	mock "game-api/pkg/domain/repository/db/user/mock"
	us "game-api/pkg/domain/service/user"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mocks struct {
	userRepository *mock.MockRepository
}

func newWithMocks(t *testing.T) (us.Service, *mocks) {
	ctrl := gomock.NewController(t)
	userRepository := mock.NewMockRepository(ctrl)
	return us.NewService(userRepository), &mocks{
		userRepository: userRepository,
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		name string
	}
	type expected struct {
		authToken string
	}
	for name, tt := range map[string]struct {
		args     args
		prepare  func(f *mocks)
		expected expected
	}{
		"【正常系】ユーザの新規登録": {
			args: args{name: "sample"},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().Create(gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(""), "sample").Return(nil).Times(1)
			},
			expected: expected{authToken: "token"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			u, m := newWithMocks(t)
			tt.prepare(m)
			got, err := u.Create(tt.args.name)
			assert.Regexp(t, regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`), got)
			assert.NoError(t, err)
		})
	}
}

func TestSelectByAuthToken(t *testing.T) {
	type args struct {
		authToken string
	}
	type expected struct {
		user *um.User
	}
	for name, tt := range map[string]struct {
		args     args
		prepare  func(f *mocks)
		expected expected
	}{
		"【正常系】ユーザの取得": {
			args: args{authToken: "authToken"},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().SelectByAuthToken("authToken").Return(&um.User{AuthToken: "authToken", Name: "sample", HighScore: 0, Coin: 0}, nil).Times(1)
			},
			expected: expected{user: &um.User{AuthToken: "authToken", Name: "sample", HighScore: 0, Coin: 0}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			u, m := newWithMocks(t)
			tt.prepare(m)
			got, err := u.SelectByAuthToken(tt.args.authToken)
			assert.Equal(t, tt.expected.user, got)
			assert.NoError(t, err)
		})
	}
}

func TestUpdateName(t *testing.T) {
	type args struct {
		user *um.User
		name string
	}
	for name, tt := range map[string]struct {
		args    args
		prepare func(f *mocks)
	}{
		"【正常系】ユーザの更新": {
			args: args{user: &um.User{ID: "id", AuthToken: "authToken", Name: "before_update", HighScore: 0, Coin: 0}, name: "updated_name"},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().Update(&um.User{ID: "id", AuthToken: "authToken", Name: "updated_name", HighScore: 0, Coin: 0}).Return(nil).Times(1)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			u, m := newWithMocks(t)
			tt.prepare(m)
			err := u.UpdateName(tt.args.user, tt.args.name)
			assert.NoError(t, err)
		})
	}
}
