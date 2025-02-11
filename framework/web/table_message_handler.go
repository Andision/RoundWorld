package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

func tableMessageHandler(ctx context.Context, t *Table, message *tableMessage) {
	sender := message.sender
	rawMessage := message.message

	parse, err := t.gameState.Parse(ctx, rawMessage)
	if err != nil {
		rawResponse, err := json.Marshal(
			tableResponse{
				TableId: t.tableId,
				Code:    tableResponseCodeError,
				Message: "parse error",
				Data:    nil,
			})
		if err != nil {
			log.Printf("Failed to marshal tableResponse: %v", err)
			return
		}
		sender.send <- rawResponse
		return
	}

	t.gameState.Lock()
	defer t.gameState.Unlock()

	valid, err := t.gameState.Validate(ctx, parse)
	if err != nil {
		rawResponse, err := json.Marshal(
			tableResponse{
				TableId: t.tableId,
				Code:    tableResponseCodeError,
				Message: fmt.Sprintf("validate error: %v", err),
				Data:    nil,
			})
		if err != nil {
			log.Printf("Failed to marshal tableResponse: %v", err)
			return
		}
		sender.send <- rawResponse
		return
	}
	if !valid {
		rawResponse, err := json.Marshal(
			tableResponse{
				TableId: t.tableId,
				Code:    tableResponseCodeError,
				Message: "invalid action",
				Data:    nil,
			})
		if err != nil {
			log.Printf("Failed to marshal tableResponse: %v", err)
			return
		}
		sender.send <- rawResponse
		return
	}

	err = t.gameState.Execute(ctx, parse)
	if err != nil {
		rawResponse, err := json.Marshal(
			tableResponse{
				TableId: t.tableId,
				Code:    tableResponseCodeError,
				Message: fmt.Sprintf("execute error: %v", err),
				Data:    nil,
			})
		if err != nil {
			log.Printf("Failed to marshal tableResponse: %v", err)
			return
		}
		sender.send <- rawResponse
		return
	}

	log.Print("execute finished")

	rawState, err := t.gameState.Encode(ctx)
	if err != nil {
		rawResponse, err := json.Marshal(
			tableResponse{
				TableId: t.tableId,
				Code:    tableResponseCodeError,
				Message: "encode error",
				Data:    nil,
			})
		if err != nil {
			log.Printf("Failed to marshal tableResponse: %v", err)
			return
		}
		sender.send <- rawResponse
		return
	}

	rawResponseExecutor, err := json.Marshal(
		tableResponse{
			TableId: t.tableId,
			Code:    tableResponseCodeSuccess,
			Message: "success",
			Data:    nil,
		},
	)
	if err != nil {
		log.Printf("Failed to marshal tableResponse: %v", err)
		return
	}
	sender.send <- rawResponseExecutor

	log.Print("broadcasting...")

	rawResponseBroadcast, err := json.Marshal(
		tableResponse{
			TableId: t.tableId,
			Code:    0,
			Message: "update",
			Data:    string(rawState),
		})
	if err != nil {
		log.Printf("Failed to marshal tableResponse: %v", err)
		return
	}
	t.broadcast <- rawResponseBroadcast
}
