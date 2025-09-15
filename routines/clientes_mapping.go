package routines

// import (
// 	"context"
// 	"fmt"
// 	"main/entities"
// 	"main/pkg/database"
// 	"main/pkg/memory"
// )

// func ClientesMemoryMapping() {

// 	consumerName := "ClientesMemoryMapping"

// 	var clients []entities.Cliente
// 	ctx := context.Background()

// 	cache := memory.GetCacheInstance()

// 	db := database.GetDB()
// 	err := db.NewSelect().Model(&clients).OrderExpr("id_client ASC").Limit(10).Scan(ctx)

// 	if err != nil {
// 		fmt.Printf("[%s] Erro ao recuperar os clients do database principal %v:\n", consumerName, err)
// 		return
// 	}

// 	// Criando um cache em memória dos amountes que não mudam
// 	// Será utilizado para verificar o limit e verificar se o cliente existe
// 	// for _, u := range clients {
// 	// 	cache.Set("cliente:"+u.ID, u.ID)
// 	// 	cache.Set("limit:"+u.ID, u.Limit)
// 	// }

// 	fmt.Printf("[%s] Mapping de memória finalizado:\n", consumerName)
// }
