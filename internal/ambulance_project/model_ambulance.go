/*
 * Meal Orders API
 *
 * Meal Orders management for Web-In-Cloud system
 *
 * API version: 1.0.0
 * Contact: <your_email>
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package ambulance_project

type Ambulance struct {

	// Unique identifier of the ambulance
	Id string `json:"id"`

	// Human readable display name of the ambulance
	Name string `json:"name"`

	RoomNumber string `json:"roomNumber"`

	MealOrders []MealOrder `json:"mealOrders,omitempty"`
}
