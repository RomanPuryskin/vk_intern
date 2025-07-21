package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vk_intern/internal/advertisements"
)

func LoadAdvertisement(ctx context.Context, adv *advertisements.CreateAdvertisementRequest) (*advertisements.Advertisement, error) {
	query := "INSERT INTO advertisements (title,description,price,image_url,login,created_at) VALUES ($1,$2,$3,$4,$5,$6)"

	created_at := time.Now()
	if _, err := Pool.Exec(ctx, query, adv.Title, adv.Description, adv.Price, adv.ImageURL, adv.UserLogin, created_at); err != nil {
		return nil, fmt.Errorf("[LoadAdvertisement|exec load advertisement]: %w", err)
	}

	return &advertisements.Advertisement{
		Title:       adv.Title,
		Description: adv.Description,
		Price:       adv.Price,
		ImageURL:    adv.ImageURL,
		UserLogin:   adv.UserLogin,
		CreatedAt:   created_at,
	}, nil
}

func GetAllAdvertisements(ctx context.Context, login string, params *advertisements.AdvertisementFilter) ([]*advertisements.AdvertisementResponse, error) {
	advs := []*advertisements.AdvertisementResponse{}

	order := params.Order
	order_by := params.OrderBy
	offset := (params.Page - 1) * params.Limit
	query := `SELECT title,description,price,image_url,login,created_at 
			FROM advertisements 
			WHERE price BETWEEN $1 AND $2
			ORDER BY ` + order_by + " " + order +
		` LIMIT $3 OFFSET $4`

	rows, err := Pool.Query(ctx, query, params.MinPrice, params.MaxPrice, params.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("[GetAllAdvertisements|exec get advs] %w", err)
	}

	for rows.Next() {
		var curAdv advertisements.AdvertisementResponse
		err := rows.Scan(&curAdv.Title, &curAdv.Description, &curAdv.Price, &curAdv.ImageURL, &curAdv.UserLogin, &curAdv.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("[GetAllAdvertisements|exec get adv] %w", err)
		}

		if curAdv.UserLogin == login {
			curAdv.IsMine = true
		}

		advs = append(advs, &curAdv)
	}
	return advs, nil
}
