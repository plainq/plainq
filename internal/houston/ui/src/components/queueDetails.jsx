import { Tabs, TabsContent } from "@/components/ui/tabs";

export default function QueueDetails({ queueDetails, error }) {
  return (
    <div className="bg-white px-2">
      <Tabs defaultValue="queue" className="w-full">
        <TabsContent value="queue">
          <div className="flex flex-row justify-between pt-4 pb-4">
            <div>
              <p className="text-2xl font-bold">Queue: {queueDetails.name}</p>
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
