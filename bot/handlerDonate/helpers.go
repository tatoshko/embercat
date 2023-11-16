package handlerDonate

import (
	"errors"
	"fmt"
	"gopkg.in/redis.v3"
	"strconv"
	"strings"
)

func getDonatesList(donates *redis.ZSliceCmd) string {
	msg := "Список донатеров:\n\n"
	for _, d := range donates.Val() {
		m := fmt.Sprintf("%s", d.Member)
		msg += fmt.Sprintf("%s %*.2f\n", m, 20-len(m), d.Score)
	}
	return fmt.Sprintf("<code>%s</code>", msg)
}

func getUserDonateMessage() string {
	return "Куда полез, пёс? Твоё дело - донатить, а награду тебе выдадут без твоего участия"
}

func getDonateMessage(info redis.Z) string {
	return fmt.Sprintf("Объявляем благодарность товарищу %s за донат в %.2f р!", info.Member, info.Score)
}

func newDonateInfo(args string) (redis.Z, error) {
	argList := strings.Split(args, " ")
	if len(argList) < 2 {
		return redis.Z{}, errors.New("incorrect number of arguments")
	}

	if val, err := strconv.ParseFloat(argList[1], 64); err != nil {
		return redis.Z{}, err
	} else {
		return redis.Z{
			Member: argList[0],
			Score:  val,
		}, nil
	}
}
