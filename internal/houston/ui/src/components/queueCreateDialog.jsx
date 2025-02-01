import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

export function QueueCreateDialog({ 
  isOpen, 
  onOpenChange,
  onCreateQueue,
  queueName,
  setQueueName,
  retentionPeriod,
  setRetentionPeriod,
  visibilityTimeout,
  setVisibilityTimeout,
  receiveAttempts,
  setReceiveAttempts,
  dropPolicy,
  setDropPolicy,
  deadLetterQueueId,
  setDeadLetterQueueId,
}) {
  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[625px]">
        <DialogHeader>
          <DialogTitle>Create Queue</DialogTitle>
          <DialogDescription>
            Enter the details for your new queue. All fields marked with * are required.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="queueName" className="text-right">
              Queue Name *
            </Label>
            <Input
              id="queueName"
              value={queueName}
              onChange={(e) => setQueueName(e.target.value)}
              className="col-span-3"
              placeholder="my-queue"
              required
            />
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="retentionPeriod" className="text-right">
              Retention Period (seconds) *
            </Label>
            <Input
              id="retentionPeriod"
              type="number"
              value={retentionPeriod}
              onChange={(e) => setRetentionPeriod(e.target.value)}
              className="col-span-3"
              min="1"
              placeholder="86400"
              required
            />
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="visibilityTimeout" className="text-right">
              Visibility Timeout (seconds) *
            </Label>
            <Input
              id="visibilityTimeout"
              type="number"
              value={visibilityTimeout}
              onChange={(e) => setVisibilityTimeout(e.target.value)}
              className="col-span-3"
              min="1"
              placeholder="30"
              required
            />
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="receiveAttempts" className="text-right">
              Max Receive Attempts *
            </Label>
            <Input
              id="receiveAttempts"
              type="number"
              value={receiveAttempts}
              onChange={(e) => setReceiveAttempts(e.target.value)}
              className="col-span-3"
              min="1"
              placeholder="5"
              required
            />
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="dropPolicy" className="text-right">
              Eviction Policy *
            </Label>
            <Select
              value={dropPolicy}
              onValueChange={setDropPolicy}
            >
              <SelectTrigger className="col-span-3">
                <SelectValue placeholder="Select a policy" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="EVICTION_POLICY_DROP">Drop Message</SelectItem>
                <SelectItem value="EVICTION_POLICY_DEAD_LETTER">Dead Letter Queue</SelectItem>
                <SelectItem value="EVICTION_POLICY_REORDER">Reorder Message</SelectItem>
              </SelectContent>
            </Select>
          </div>
          {dropPolicy === "EVICTION_POLICY_DEAD_LETTER" && (
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="deadLetterQueueId" className="text-right">
                Dead Letter Queue ID
              </Label>
              <Input
                id="deadLetterQueueId"
                value={deadLetterQueueId}
                onChange={(e) => setDeadLetterQueueId(e.target.value)}
                className="col-span-3"
                placeholder="Enter queue ID"
              />
            </div>
          )}
        </div>
        <DialogFooter>
          <Button type="submit" onClick={onCreateQueue}>
            Create Queue
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
