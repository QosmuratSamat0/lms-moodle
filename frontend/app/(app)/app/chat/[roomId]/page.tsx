import ChatClient from "./ChatClient";

export async function generateStaticParams() {
  // Return empty array - room IDs are dynamic and fetched at runtime
  // Pages will be rendered on-demand (ISR/dynamic rendering)
  return [];
}

export default function ChatRoomPage({ params }: { params: { roomId: string } }) {
  return <ChatClient roomId={params.roomId} />;
}
