package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nonamecat19/go-orm/orm/lib/querybuilder"
	"shopping-list/backend/database"
	"shopping-list/backend/entities"
	"strconv"
)

type ItemCreate struct {
	Name string `json:"name"`
}

func CreateItem(c *fiber.Ctx) error {
	var itemCreate ItemCreate
	if err := c.BodyParser(&itemCreate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Непрвильний тип даних",
		})
	}

	if itemCreate.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Назва товару обов'язкова",
		})
	}

	newItem := entities.Item{
		Name:   itemCreate.Name,
		Bought: false,
	}

	err := querybuilder.
		CreateQueryBuilder(database.DbClient).
		InsertOne(newItem)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Не вдалось додати товар: %s", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Товар створено успішно",
	})
}

func GetItems(c *fiber.Ctx) error {
	var items []entities.Item
	err := querybuilder.CreateQueryBuilder(database.DbClient).
		Select("items.id", "items.name", "items.bought").
		FindMany(&items)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Не вдалось отримати товари: %s", err),
		})
	}

	return c.JSON(items)
}

type ItemUpdate struct {
	Name   string `json:"name,omitempty"`
	Bought bool   `json:"bought,omitempty"`
}

func UpdateItem(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid item ID",
		})
	}

	var itemUpdate ItemUpdate
	if err = c.BodyParser(&itemUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = querybuilder.CreateQueryBuilder(database.DbClient).
		Where("id = ?", id).
		SetValues(map[string]any{"name": itemUpdate.Name, "bought": itemUpdate.Bought}).
		UpdateMany(&entities.Item{})

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update item",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Товар оновлено успішно",
	})
}

func DeleteItem(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid item ID",
		})
	}

	err = querybuilder.CreateQueryBuilder(database.DbClient).
		Where("id = ?", id).
		DeleteMany(&entities.Item{})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete item",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Товар видалено успішно",
	})
}
