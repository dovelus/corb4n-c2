import { useState } from "react";
import { useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import {
  Clock,
  Cpu,
  HardDrive,
  Download,
  Terminal,
  Camera,
  X,
  File,
  Copy,
  Check,
  Eye,
} from "lucide-react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "../components/ui/alert-dialog";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../components/ui/dialog";
import {
  fetchBeaconById,
  fetchBeaconFiles,
  Beacon,
  BeaconFile,
} from "../api/beacons";

export default function BeaconDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const [command, setCommand] = useState("");
  const [previewFile, setPreviewFile] = useState<BeaconFile | null>(null);
  const [isCopied, setIsCopied] = useState(false);
  const [isPreviewOpen, setIsPreviewOpen] = useState(false);

  const {
    data: beaconData,
    isLoading: isLoadingBeacon,
    error: beaconError,
  } = useQuery<Beacon, Error>({
    queryKey: ["beacon", id],
    queryFn: () => fetchBeaconById(id!),
    enabled: !!id,
  });

  const {
    data: beaconFiles,
    isLoading: isLoadingFiles,
    error: filesError,
  } = useQuery<BeaconFile[], Error>({
    queryKey: ["beaconFiles", id],
    queryFn: () => fetchBeaconFiles(id!),
    enabled: !!id,
  });

  if (isLoadingBeacon) return <div>Loading beacon details...</div>;
  if (beaconError)
    return (
      <div>
        An error occurred while fetching beacon details: {beaconError.message}
      </div>
    );
  if (!beaconData) return <div>No beacon data found</div>;

  const handleDownload = (fileIdentifier: string, beaconId: string) => {
    console.log(`Downloading file: ${fileIdentifier} for beacon: ${beaconId}`);
    // Implement actual download logic here
  };

  const handlePreview = (file: BeaconFile) => {
    setPreviewFile(file);
    setIsPreviewOpen(true);
  };

  const handleCopyContent = () => {
    if (previewFile) {
      navigator.clipboard
        .writeText(previewFile.Output)
        .then(() => {
          setIsCopied(true);
          setTimeout(() => setIsCopied(false), 2000); // Reset after 2 seconds
        })
        .catch((err) => {
          console.error("Failed to copy content: ", err);
        });
    }
  };

  const executeCommand = () => {
    console.log(`Executing command: ${command}`);
    setCommand("");
  };

  const captureScreenshot = () => {
    console.log("Capturing screenshot");
  };

  const killBeacon = () => {
    console.log("Killing beacon");
  };

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">
        Beacon Details: {beaconData.beaconID}
      </h1>
      <div className="grid gap-4 md:grid-cols-5">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Operating System
            </CardTitle>
            <HardDrive className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{beaconData.os}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Running User</CardTitle>
            <HardDrive className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{beaconData.processUser}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Endpoint Protection
            </CardTitle>
            <HardDrive className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{beaconData.endpointProt}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Last Checkout</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {new Intl.DateTimeFormat("en-GB", {
                day: "2-digit",
                month: "2-digit",
                year: "2-digit",
                hour: "2-digit",
                minute: "2-digit",
                second: "2-digit",
                hour12: false,
              }).format(new Date(beaconData.lastUpdate))}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">PID</CardTitle>
            <Cpu className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{beaconData.processID}</div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle>Execute Command</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex space-x-2">
              <Input
                type="text"
                placeholder="Enter command"
                value={command}
                onChange={(e) => setCommand(e.target.value)}
              />
              <Button onClick={executeCommand}>
                <Terminal className="h-4 w-4 mr-2" />
                Execute
              </Button>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Capture Screenshot</CardTitle>
          </CardHeader>
          <CardContent>
            <Button onClick={captureScreenshot} className="w-full">
              <Camera className="h-4 w-4 mr-2" />
              Capture Screenshot
            </Button>
          </CardContent>
        </Card>
        <Card className="bg-red-500 dark:bg-red-900">
          <CardHeader>
            <CardTitle className="text-white">Kill Beacon</CardTitle>
          </CardHeader>
          <CardContent>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button className="w-full bg-white text-red-500 hover:bg-red-100">
                  <X className="h-4 w-4 mr-2" />
                  Kill Beacon
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                  <AlertDialogDescription>
                    This action cannot be undone. This will permanently kill the
                    beacon process.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction
                    onClick={killBeacon}
                    className="bg-red-500 hover:bg-red-600 text-white"
                  >
                    Confirm Kill
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Beacon Files</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoadingFiles ? (
            <div className="text-center py-4 text-muted-foreground">
              Loading files...
            </div>
          ) : filesError ? (
            <div className="text-center py-4 text-red-500">
              Error loading files: {filesError.message}
            </div>
          ) : beaconFiles && beaconFiles.length > 0 ? (
            <ul className="space-y-2">
              {beaconFiles.map((file, index) => (
                <li
                  key={index}
                  className="flex items-center justify-between p-2 bg-muted rounded-md"
                >
                  <div className="flex items-center">
                    <File className="h-4 w-4 mr-2" />
                    <span className="font-medium">{file.FileName}</span>
                    <span className="text-sm text-muted-foreground ml-2">
                      ({file.FileType})
                    </span>
                  </div>
                  <div className="flex space-x-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handlePreview(file)}
                      className="text-[#78716C] dark:text-[#A8A29E] border-[#E8E8E8] dark:border-[#292524] hover:bg-[#F5F5F4] dark:hover:bg-[#1C1917]"
                    >
                      <Eye className="h-4 w-4 mr-2" />
                      Preview
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() =>
                        handleDownload(file.FileName, file.BeaconID)
                      }
                      className="text-[#78716C] dark:text-[#A8A29E] border-[#E8E8E8] dark:border-[#292524] hover:bg-[#F5F5F4] dark:hover:bg-[#1C1917]"
                    >
                      <Download className="h-4 w-4 mr-2" />
                      Download
                    </Button>
                  </div>
                </li>
              ))}
            </ul>
          ) : (
            <div className="text-center py-4 text-muted-foreground">
              No files available for this beacon
            </div>
          )}
        </CardContent>
      </Card>

      <Dialog
        open={isPreviewOpen}
        onOpenChange={(open) => {
          setIsPreviewOpen(open);
          if (!open) setPreviewFile(null);
        }}
      >
        <DialogContent className="sm:max-w-[800px] max-h-[80vh] overflow-hidden flex flex-col">
          {previewFile && (
            <div className="pt-4 flex flex-col h-full">
              <DialogHeader className="flex-shrink-0 flex flex-col sm:flex-row justify-between items-start sm:items-center space-y-2 sm:space-y-0">
                <DialogTitle className="text-left text-lg font-semibold">
                  File Preview: {previewFile.FileName}
                </DialogTitle>
                <Button
                  variant="outline"
                  onClick={handleCopyContent}
                  className="relative w-full sm:w-auto"
                >
                  <span
                    className={`flex items-center justify-center transition-opacity duration-300 ${isCopied ? "opacity-0" : "opacity-100"}`}
                  >
                    <Copy className="h-4 w-4 mr-2" />
                    Copy Content
                  </span>
                  <span
                    className={`absolute inset-0 flex items-center justify-center transition-opacity duration-300 ${isCopied ? "opacity-100" : "opacity-0"}`}
                  >
                    <Check className="h-4 w-4 mr-2 text-green-500" />
                    Copied!
                  </span>
                </Button>
              </DialogHeader>
              <div className="mt-4 flex-grow overflow-hidden">
                <pre className="h-full overflow-auto rounded-md bg-muted p-4 text-sm">
                  <code className="inline-block min-w-full whitespace-pre">
                    {previewFile.Output.split("\n").map((line, index) => (
                      <div key={index} className="flex">
                        <span className="select-none inline-block w-12 mr-4 text-right text-muted-foreground">
                          {index + 1}
                        </span>
                        <span>{line}</span>
                      </div>
                    ))}
                  </code>
                </pre>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
