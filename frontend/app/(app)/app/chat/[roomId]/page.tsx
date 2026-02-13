import ChatClient from "./ChatClient";

export default function ChatRoomPage({ params }: { params: { roomId: string } }) {
  return <ChatClient roomId={params.roomId} />;
}
