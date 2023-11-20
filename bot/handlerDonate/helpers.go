package handlerDonate

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
)

func getDonatesList(donates Donates) string {
    msg := "Список донатеров:\n\n"
    for _, d := range donates {
        m := fmt.Sprintf("%s", d.Username)
        msg += fmt.Sprintf("%s %*.2f\n", m, 20-len(m), d.Sum)
    }
    return fmt.Sprintf("<code>%s</code>", msg)
}

func getUserDonateMessage() string {
    return "Куда полез, пёс? Твоё дело - донатить, а награду тебе выдадут без твоего участия"
}

func getDonateMessage(donater string, score float64) string {
    return fmt.Sprintf("Объявляем благодарность товарищу %s за донат в %.2f р!", donater, score)
}

func parseArgs(args string) (donater string, sum float64, err error) {
    argList := strings.Split(args, " ")

    if len(argList) < 2 {
        return "", 0.0, errors.New("incorrect number of arguments, must be 2")
    }

    if sum, err = strconv.ParseFloat(argList[1], 64); err != nil {
        return
    }

    return argList[0], sum, nil
}
