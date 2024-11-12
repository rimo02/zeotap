package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rimo02/zeotap/assignment1/database"
)

func CreateRule(c *gin.Context) {
	var body struct {
		RuleID     string `json:"rule_id"`
		RuleString string `json:"rule_string"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Printf("The rule string is: %s", body.RuleString)

	ast, err := parseRuleString(body.RuleString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule string"})
		return
	}

	newRule := &database.Rule{
		RuleString: body.RuleString,
		AST:        ast,
		RuleID:     body.RuleID,
	}
	collection := database.GetCollection(database.Client, "rule")
	res, err := collection.InsertOne(context.TODO(), newRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert rule into the database"})
		return
	}
	fmt.Printf("Inserted rule with ID: %v\n", res.InsertedID)

	c.JSON(http.StatusCreated, newRule)
}
