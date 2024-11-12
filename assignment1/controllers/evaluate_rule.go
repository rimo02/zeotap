package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rimo02/zeotap/assignment1/database"
	"github.com/rimo02/zeotap/assignment1/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

func compare(left, right interface{}, operator string) (bool, error) {
	switch leftValue := left.(type) {
	case float64:
		// Ensure right is also a float64
		rightValue, ok := right.(float64)
		if !ok {
			return false, fmt.Errorf("type mismatch: %v and %v", leftValue, right)
		}
		// Perform numeric comparisons
		switch operator {
		case ">":
			return leftValue > rightValue, nil
		case ">=":
			return leftValue >= rightValue, nil
		case "<":
			return leftValue < rightValue, nil
		case "<=":
			return leftValue <= rightValue, nil
		default:
			return false, fmt.Errorf("unknown operator: %s", operator)
		}

	case int: // Handle integer comparisons as well
		rightValue, ok := right.(int)
		if !ok {
			return false, fmt.Errorf("type mismatch: %v and %v", leftValue, right)
		}
		switch operator {
		case ">":
			return leftValue > rightValue, nil
		case ">=":
			return leftValue >= rightValue, nil
		case "<":
			return leftValue < rightValue, nil
		case "<=":
			return leftValue <= rightValue, nil
		default:
			return false, fmt.Errorf("unknown operator: %s", operator)
		}

	case string:
		// String comparisons (optional, can extend to >, < for lexicographical order)
		rightValue, ok := right.(string)
		if !ok {
			return false, fmt.Errorf("type mismatch: %v and %v", leftValue, right)
		}
		switch operator {
		case "==":
			return leftValue == rightValue, nil
		case "!=":
			return leftValue != rightValue, nil
		default:
			return false, fmt.Errorf("unsupported operator for strings: %s", operator)
		}

	default:
		return false, fmt.Errorf("unsupported type: %T", left)
	}
}

func evaluateNode(node *model.Node, data map[string]interface{}) (bool, error) {
	// Base case: if the node is nil, return false with no error
	if node == nil {
		return false, nil
	}

	// If the node is an operator node (AND, OR)
	if node.Type == "operator" {
		attr := node.Value.(string)
		left, errL := evaluateNode(node.Left, data)
		if errL != nil {
			return false, errL
		}
		right, errR := evaluateNode(node.Right, data)
		if errR != nil {
			return false, errR
		}
		// Handling logical operations
		if attr == "AND" {
			return left && right, nil // Both sides must be true
		}
		if attr == "OR" {
			return left || right, nil // At least one side must be true
		}
		return false, fmt.Errorf("unknown operator: %s", attr)
	}

	// Operand node: evaluate condition
	if node.Type == "operand" {
		val := node.Value.(map[string]interface{})
		attr := val["attribute"].(string)
		operator := val["operator"].(string)
		value := val["value"]

		fmt.Println("Val = ", val)
		fmt.Println("Attr = ", attr)
		fmt.Println("operator = ", operator)
		fmt.Println("value = ", value)

		// Retrieve the attribute value from the data map
		dataValue, exists := data[attr]
		if !exists {
			return false, fmt.Errorf("attribute %s not found in data", attr)
		}

		// Handle comparison based on the operator and value types
		switch operator {
		case "=": // Handle equality check for strings
			return dataValue == value, nil
		case ">":
			return compare(dataValue, value, operator)
		case ">=":
			return compare(dataValue, value, operator)
		case "<":
			return compare(dataValue, value, operator)
		case "<=":
			return compare(dataValue, value, operator)
		case "!=":
			return dataValue != value, nil // Handle inequality check
		default:
			return false, fmt.Errorf("unknown operator: %s", operator)
		}
	}

	return false, fmt.Errorf("unknown node type: %s", node.Type)
}

func EvaluateRule(c *gin.Context) {
	var body struct {
		RuleID string                 `json:"ruleId"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var result []string
	for key, value := range body.Data {
		result = append(result, fmt.Sprintf("%s = %v", key, value))
	}

	output := strings.Join(result, " AND ")
	fmt.Println(output)

	var ast *model.Node
	collection := database.GetCollection(database.Client, "rule")
	err := collection.FindOne(context.TODO(), bson.M{"ruleid": body.RuleID}).Decode(&ast)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"message": "rule not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	evaluationResult, err := evaluateNode(ast, body.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": evaluationResult})
}
