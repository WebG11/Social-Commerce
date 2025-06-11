namespace go agent

service AgentService {
    AskResp ask(1: AskReq req);

    ChatMessageResp SendMessage(1: ChatMessageReq req)

    HistoryMessageResp GetHistoryMessages(1: HistoryMessageReq req)
}


struct AskReq {
    1: string user_prompt;
    2: i32 user_id;
}

struct AskResp {
    1: string content;
}

struct ChatMessageReq {
    1: required i64 sender_id,
    2: required i64 receiver_id,
    3: required string content,
}

struct ChatMessageResp {
}

struct HistoryMessageReq {
    1: required i64 sender_id,
}

struct HistoryMessageResp {
    1: required list<ChatMessage> messages,
}

struct ChatMessage {
    1: required i64 id,
    2: required i64 sender_id,
    3: required i64 receiver_id,
    4: required string content,
    5: required string timestamp,
}