package order

//func Test_GetOrderByID(t *testing.T) {
//	t.Parallel()
//
//	t.Run("success", func(t *testing.T) {
//		t.Parallel()
//
//		//arrange
//		ctx := context.Background()
//		mc := minimock.NewController(t)
//		repoMock := mock.NewRepoMock(mc).GetOrderByIDMock.Expect(ctx, fixtures.ID).Return(fixtures.PVZ().V(), nil)
//		service := New(repoMock)
//
//		// act
//		order, err := service.GetOrderByID(ctx, fixtures.ID)
//
//		// assert
//		require.NoError(t, err)
//		assert.Equal(t, order, fixtures.PVZ().V())
//	})
//
//	t.Run("not found", func(t *testing.T) {
//		t.Parallel()
//
//		//arrange
//		ctx := context.Background()
//		mc := minimock.NewController(t)
//		repoMock := mock.NewRepoMock(mc).GetOrderByIDMock.Expect(ctx, fixtures.ID).Return(model.Order{}, repo2.ErrNotFound)
//		service := New(repoMock)
//
//		// act
//		order, err := service.GetOrderByID(ctx, fixtures.ID)
//
//		// assert
//		require.Equal(t, err.Error(), "s.repo.GetOrderByID: not found")
//		assert.Equal(t, order, model.Order{})
//	})
//
//	t.Run("db.Get error", func(t *testing.T) {
//		t.Parallel()
//
//		//arrange
//		ctx := context.Background()
//		mc := minimock.NewController(t)
//		repoMock := mock.NewRepoMock(mc).GetOrderByIDMock.Expect(ctx, fixtures.ID).Return(model.Order{}, fmt.Errorf("repo.db.Get: %s", "error"))
//		service := New(repoMock)
//
//		// act
//		order, err := service.GetOrderByID(ctx, fixtures.ID)
//
//		// assert
//		require.Equal(t, err.Error(), "s.repo.GetOrderByID: repo.db.Get: error")
//		assert.Equal(t, order, model.Order{})
//	})
//
//	tests := []struct {
//		name        string
//		wantOrder   model.Order
//		returnOrder model.Order
//		returnErr   error
//		wantErr     error
//	}{
//		{
//			name:        "success",
//			wantOrder:   fixtures.PVZ().V(),
//			returnOrder: fixtures.PVZ().V(),
//			returnErr:   nil,
//			wantErr:     nil,
//		},
//		{
//			name:        "not found",
//			returnOrder: model.Order{},
//			wantOrder:   model.Order{},
//			returnErr:   repo2.ErrNotFound,
//			wantErr:     fmt.Errorf("s.repo.GetOrderByID: %w", repo2.ErrNotFound),
//		},
//		{
//			name:        "db.Get error",
//			returnOrder: model.Order{},
//			wantOrder:   model.Order{},
//			returnErr:   fmt.Errorf("repo.db.Get: %s", "error"),
//			wantErr:     fmt.Errorf("s.repo.GetOrderByID: repo.db.Get: %s", "error"),
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			t.Parallel()
//
//			// arrange
//			ctx := context.Background()
//			mc := minimock.NewController(t)
//			repoMock := mock.NewRepoMock(mc).GetOrderByIDMock.Expect(ctx, fixtures.ID).Return(tt.returnOrder, tt.returnErr)
//			service := New(repoMock)
//
//			// act
//			order, err := service.GetOrderByID(ctx, fixtures.ID)
//
//			// assert
//			require.ErrorIs(t, err, tt.wantErr)
//			require.Equal(t, order, tt.wantOrder)
//		})
//	}
//}
