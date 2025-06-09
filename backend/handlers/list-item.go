package handlers

import (
	"fmt"
	"github.com/nonamecat19/go-orm/orm/lib/querybuilder"
	"shopping-list/backend/database"
	"shopping-list/backend/entities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AddItemToListBody struct {
	ItemId int `json:"itemId"`
	ListId int `json:"listId"`
}

func AddItemToList(c *fiber.Ctx) error {
	var addItemToList AddItemToListBody

	if err := c.BodyParser(&addItemToList); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	fmt.Println(addItemToList)

	err := querybuilder.CreateQueryBuilder(database.DbClient).
		Where("id = ?", addItemToList.ItemId).
		SetValues(map[string]any{"list_id": addItemToList.ListId}).
		Debug().
		UpdateMany(&entities.Item{})

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Помилка додавання продукту до списку",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Продукт доданий у список успішно",
	})
}

func RemoveItemFromList(c *fiber.Ctx) error {
	listID, err := strconv.ParseInt(c.Params("listId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Непрвильний тип даних",
		})
	}
	itemID, err := strconv.ParseInt(c.Params("itemId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Непрвильний тип даних",
		})
	}

	err = querybuilder.CreateQueryBuilder(database.DbClient).
		Where("id = ?", itemID).
		AndWhere("list_id = ?", listID).
		SetValues(map[string]any{"list_id": nil}).
		UpdateMany(&entities.Item{})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Помилка видалення продукту зі списку",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Продукт видалений зі списку успішно",
	})
}
