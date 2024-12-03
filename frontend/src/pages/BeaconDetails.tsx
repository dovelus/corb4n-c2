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
import { fetchBeaconById, Beacon } from "../api/beacons";

interface BeaconFiles {
  FileName: string;
  FileType: string;
  Output: string;
  BeaconID: string;
}

const mockBeaconFiles: BeaconFiles[] = [
  {
    FileName: "screenshot.png",
    FileType: "Image",
    Output: "Base64 encoded data",
    BeaconID: "beacon123",
  },
  {
    FileName: "system_info.txt",
    FileType: "Text",
    Output: "System information data",
    BeaconID: "beacon123",
  },
];

export default function BeaconDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const [command, setCommand] = useState("");
  const {
    data: beaconData,
    isLoading,
    error,
  } = useQuery<Beacon, Error>({
    queryKey: ["beacon", id],
    queryFn: () => fetchBeaconById(id!),
    enabled: !!id,
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>An error occurred: {error.message}</div>;
  if (!beaconData) return <div>No beacon data found</div>;

  const handleDownload = (fileName: string, beaconId: string) => {
    console.log(`Downloading file: ${fileName} for beacon: ${beaconId}`);
    // Implement actual download logic here
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
          <CardTitle>Downloadable Items</CardTitle>
        </CardHeader>
        <CardContent>
          <ul className="space-y-2">
            {mockBeaconFiles.map((item, index) => (
              <li
                key={index}
                className="flex items-center justify-between p-2 bg-muted rounded-md"
              >
                <div>
                  <span className="font-medium">{item.FileName}</span>
                  <span className="text-sm text-muted-foreground ml-2">
                    ({item.FileType})
                  </span>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleDownload(item.FileName, item.BeaconID)}
                  className="text-[#78716C] dark:text-[#A8A29E] border-[#E8E8E8] dark:border-[#292524] hover:bg-[#F5F5F4] dark:hover:bg-[#1C1917]"
                >
                  <Download className="h-4 w-4 mr-2" />
                  Download
                </Button>
              </li>
            ))}
          </ul>
        </CardContent>
      </Card>
    </div>
  );
}
