package client

// type BudgetClient struct {
// 	client pb.BudgetServiceClient
// }

// // NewBudgetClient creates a new gRPC client for the budget service
// func NewBudgetClient(serviceAddress string) (*BudgetClient, error) {
// 	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &BudgetClient{
// 		client: pb.NewBudgetServiceClient(conn),
// 	}, nil
// }

// // Implement gRPC client methods, e.g., CreateUser
// func (bc *BudgetClient) CreateUser(ctx context.Context, name string) (string, error) {
// 	req := &pb.CreateUserRequest{Name: name}
// 	res, err := bc.client.CreateUser(ctx, req)
// 	if err != nil {
// 		return "", err
// 	}
// 	return res.Id, nil
// }
