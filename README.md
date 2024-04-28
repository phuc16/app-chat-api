# app-chat-api

Websocket note
{
    "type": "",
    "list_user": ["newChat"],
    "chat": {
        "from": "this_is_userID",
        "to": "this_is_conversationID",
        "message": "this is message"
    }
}
if create new conversation => type = "newChat" and list_user must contain list userID in new conversation
else type = "" and list_user = []