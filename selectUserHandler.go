package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var selectedUsersMap = make(map[string][]string)

func selectUserHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 길드 멤버 가져오기
	members, err := s.GuildMembers(*GuildID, "", 25)
	if err != nil {
		log.Println("Error fetching members:", err)
		return
	}

	// User 목록으로 SelectMenu 생성
	var options []discordgo.SelectMenuOption
	for _, member := range members {
		// member.User.ID와 member.User.Username을 사용하여 옵션 생성
		if !member.User.Bot {
			options = append(options, discordgo.SelectMenuOption{
				Label: member.User.GlobalName,
				Value: member.User.ID,
			})
		}
	}
	MinValues := 1
	MaxValues := len(options)

	// SelectMenu와 ActionRow 설정
	selectMenu := discordgo.SelectMenu{
		CustomID:    "user_select_menu",
		Placeholder: "Select a user...",
		MinValues:   &MinValues,
		MaxValues:   MaxValues,
		Options:     options,
	}
	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			selectMenu,
		},
	}

	// start_button
	buttonRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			&discordgo.Button{
				Label:    "Select",                // 버튼 텍스트
				Style:    discordgo.PrimaryButton, // 버튼 스타일
				CustomID: "start_button",          // 버튼 클릭 시 처리할 ID
			},
		},
	}

	//select_all_button

	// 드롭다운 메뉴와 버튼을 포함한 메시지 전송
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				actionRow,
				buttonRow,
			},
		},
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
		return
	}

	selectedUsersMap[i.GuildID] = make([]string, 0)
}

// Select 버튼이 눌렸을 때 선택된 멤버들을 처리하는 핸들러
func handleStartButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 선택된 멤버 ID 목록을 가져옴
	tempSelectedMembers := selectedUsersMap[i.GuildID]
	if len(tempSelectedMembers) == 0 {
		log.Println("No members selected.")
		return
	}

	king := Get(tempSelectedMembers)
	var message string
	for _, v := range tempSelectedMembers {
		if v == king {
			message = "당신은 왕 입니다!"
		} else {
			message = "당신은 왕이 아닙니다!"
		}
		sendPrivateMessage(s, v, message)
	}

	// Interaction 응답 (선택 결과를 유저에게 표시)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
		// Type: discordgo.InteractionResponseChannelMessageWithSource,
		// Data: &discordgo.InteractionResponseData{
		// 	Content: content,
		// },
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
	}
}

func handleSelectMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Map 변수
	selectedUsersMap[i.GuildID] = i.MessageComponentData().Values

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 상호작용 지연
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println("Error responding to select menu interaction:", err)
	}
}
