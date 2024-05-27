package ambulance_project

import (
  "net/http"
  "log"
  "github.com/gin-gonic/gin"
  "github.com/google/uuid"
  "slices"
)

// CreateMealOrder - Creates a new meal order
func (this *implMealOrdersAPI) CreateMealOrder(ctx *gin.Context) {
    updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
        var order MealOrder

        if err := c.ShouldBindJSON(&order); err != nil {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Invalid request body",
                "error":   err.Error(),
            }, http.StatusBadRequest
        }

        log.Printf("Received order: %+v", order)

        if order.Name == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Patient name is required",
            }, http.StatusBadRequest
        }

        if order.DietaryReq == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Dietary requirements are required",
            }, http.StatusBadRequest
        }

        if order.MedicalNeed == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Medical needs are required",
            }, http.StatusBadRequest
        }

        if order.Id == "" || order.Id == "@new" {
            order.Id = uuid.NewString()
        }

        ambulance.MealOrders = append(ambulance.MealOrders, order)
        
        orderIndex := slices.IndexFunc(ambulance.MealOrders, func(o MealOrder) bool {
            return order.Id == o.Id
        })
        
        if orderIndex < 0 {
            return nil, gin.H{
                "status":  http.StatusInternalServerError,
                "message": "Failed to save order",
            }, http.StatusInternalServerError
        }
        return ambulance, ambulance.MealOrders[orderIndex], http.StatusOK
    })
}

// DeleteMealOrder - Deletes a specific meal order
func (this *implMealOrdersAPI) DeleteMealOrder(ctx *gin.Context) {
    updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
        orderId := ctx.Param("orderId")

        if orderId == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Order ID is required",
            }, http.StatusBadRequest
        }

        orderIndex := slices.IndexFunc(ambulance.MealOrders, func(order MealOrder) bool {
            return orderId == order.Id
        })

        if orderIndex < 0 {
            return nil, gin.H{
                "status":  http.StatusNotFound,
                "message": "Order not found",
            }, http.StatusNotFound
        }

        ambulance.MealOrders = append(ambulance.MealOrders[:orderIndex], ambulance.MealOrders[orderIndex+1:]...)
        return ambulance, nil, http.StatusNoContent
    })
}

// GetMealOrder - Provides details about a specific meal order
func (this *implMealOrdersAPI) GetMealOrder(ctx *gin.Context) {
    updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
        orderId := ctx.Param("orderId")

        if orderId == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Order ID is required",
            }, http.StatusBadRequest
        }

        orderIndex := slices.IndexFunc(ambulance.MealOrders, func(order MealOrder) bool {
            return orderId == order.Id
        })

        if orderIndex < 0 {
            return nil, gin.H{
                "status":  http.StatusNotFound,
                "message": "Order not found",
            }, http.StatusNotFound
        }

        return nil, ambulance.MealOrders[orderIndex], http.StatusOK
    })
}

// GetMealOrders - Provides the list of meal orders
func (this *implMealOrdersAPI) GetMealOrders(ctx *gin.Context) {
    updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
        result := ambulance.MealOrders
        if result == nil {
            result = []MealOrder{}
        }
        // return nil ambulance - no need to update it in db
        return nil, result, http.StatusOK
    })
}

// UpdateMealOrder - Updates a specific meal order
func (this *implMealOrdersAPI) UpdateMealOrder(ctx *gin.Context) {
    updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
        var updatedOrder MealOrder

        if err := c.ShouldBindJSON(&updatedOrder); err != nil {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Invalid request body",
                "error":   err.Error(),
            }, http.StatusBadRequest
        }

        orderId := ctx.Param("orderId")

        orderIndex := slices.IndexFunc(ambulance.MealOrders, func(order MealOrder) bool {
            return orderId == order.Id
        })

        if orderIndex < 0 {
            return nil, gin.H{
                "status":  http.StatusNotFound,
                "message": "Order not found",
            }, http.StatusNotFound
        }

        ambulance.MealOrders[orderIndex] = updatedOrder

        return ambulance, ambulance.MealOrders[orderIndex], http.StatusOK
    })
}
