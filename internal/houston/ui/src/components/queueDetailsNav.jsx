import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/components/ui/tabs"
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar"
import QueueDetails from "@/components/queueDetails.jsx";

export default function QueueDetailsNav({queueDetails, error}) {
  return (
    <div className="">
      <div className="flex flex-row justify-between items-center px-6">
        <div>
          <a href="/" className="text-lg text-gray-500 hover:text-gray-900">PlainQ</a>
        </div>
        <div className="flex flex-row items-center gap-2">
          <div>
            <a href="https://docs.plainq.com">Docs</a>
          </div>
          <Avatar>
            <AvatarImage src="https://github.com/shadcn.png"/>
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
        </div>
      </div>
      <div className="bg-white p-4 rounded">
        <Tabs defaultValue="queue">
          <TabsList className="w-full justify-start">
            <TabsTrigger value="queue"><span className="">Queues</span></TabsTrigger>
            <TabsTrigger value="pubsub"><span className="">PubSub</span></TabsTrigger>
            <TabsTrigger value="users&access"><span className="">Users & Access</span></TabsTrigger>
            <TabsTrigger value="settings"><span className="">Settings</span></TabsTrigger>
          </TabsList>
          <QueueDetails queueDetails={queueDetails} error={error}/>
          <TabsContent value="pubsub" className="px-1">Make changes to your account here.</TabsContent>
          <TabsContent value="users&access" className="px-1">Change your password here.</TabsContent>
          <TabsContent value="settings" className="px-1">Change your password here.</TabsContent>
        </Tabs>
      </div>
    </div>
  )
}