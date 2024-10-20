package controllers

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/rimo02/zeotap/assignment1/model"
	"context"
	"github.com/rimo02/zeotap/assignment1/database"
)

func CombineRules(c *gin.Context) {
	var body struct {
		RuleId string `json:"rule_id"`
		RuleStrings     []string `json:"ruleStrings"`
		CombineOperator string   `json:"combop"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var combinedAST *model.Node
	for i, ruleString := range body.RuleStrings {
		ast, err := parseRuleString(ruleString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing rule"})
			return
		}
		if i == 0 {
			combinedAST = ast
		} else {
			combinedAST = &model.Node{
				Type:  "operator",
				Left:  combinedAST,
				Right: ast,
				Value: strings.ToUpper(body.CombineOperator),
			}
		}
	}
	newRule := &database.Rule{
		RuleString: strings.Join(body.RuleStrings,""),
		AST:        combinedAST,
		RuleID:     body.RuleId,
	}
	collection := database.GetCollection(database.Client, "rule")
	collection.InsertOne(context.TODO(), newRule)

	c.JSON(http.StatusOK, gin.H{"combinedAST": combinedAST})
}
