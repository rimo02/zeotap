package controllers

import (
	"errors"
	"fmt"
	"github.com/rimo02/zeotap/assignment1/model"
	"strings"
)

// Function to parse rule string into an AST
func parseRuleString(ruleString string) (*model.Node, error) {
	fmt.Println("Parsing rule string:", ruleString)
	tokens := strings.Fields(ruleString)
	if len(tokens) == 0 {
		return nil, errors.New("invalid rule string")
	}

	outputQueue := []string{}
	operatorStack := []string{}
	precedence := map[string]int{"AND": 1, "OR": 0}

	// conver the given expression to postfix expression
	// sum > 30 AND department = sales
	// sum > 30  department  = sales AND --- postfix
	for _, token := range tokens {
		switch token {
		case "AND", "OR":
			for len(operatorStack) > 0 && precedence[operatorStack[len(operatorStack)-1]] >= precedence[token] {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, token)
		// case "(":
		// 	operatorStack = append(operatorStack, token)
		// case ")":
		// 	for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
		// 		outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
		// 		operatorStack = operatorStack[:len(operatorStack)-1]
		// 	}
		// 	operatorStack = operatorStack[:len(operatorStack)-1]
		case "(", ")":
			continue
		default:
			outputQueue = append(outputQueue, token)
		}
	}
	for len(operatorStack) > 0 {
		outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	fmt.Println("Output queue = ",outputQueue)

	var stack []*model.Node

	createOperandNode := func(attribute, operator, value string) *model.Node {
		return &model.Node{
			Type:  "operand",
			Value: map[string]interface{}{"attribute": attribute, "operator": operator, "value": value},
		}
	}

	for i := 0; i < len(outputQueue); i++ {
		token := outputQueue[i]
		if token == "AND" || token == "OR" {
			right := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			left := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, &model.Node{Type: "operator", Left: left, Right: right, Value: token})
		} else if strings.ContainsAny(token, "><=") {
			value := outputQueue[i+1]
			i++
			attribute := stack[len(stack)-1].Value.(string)
			stack = stack[:len(stack)-1]
			stack = append(stack, createOperandNode(attribute, token, value))
		} else {
			stack = append(stack, &model.Node{Type: "operand", Value: token})
		}
	}

	if len(stack) != 1 {
		return nil, errors.New("invalid rule string")
	}

	return stack[0], nil
}
