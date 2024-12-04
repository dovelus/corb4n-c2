import axios from "axios";

export interface Beacon {
  beaconID: string;
  hostname: string;
  beaconIntIP: string;
  beaconExtIP: string;
  os: string;
  endpointProt: string;
  processID: number;
  processUser: string;
  lastUpdate: string;
}

export interface BeaconFile {
  FileName: string;
  FileType: string;
  Output: string;
  BeaconID: string;
}

export const fetchBeacons = async (): Promise<Beacon[]> => {
  const response = await axios.get<Beacon[]>("/api/beacons");
  return response.data;
};

export const fetchBeaconById = async (id: string): Promise<Beacon> => {
  const response = await axios.get<Beacon[]>("/api/beacons");
  const beacon = response.data.find((b) => b.beaconID === id);
  if (!beacon) {
    throw new Error(`Beacon with id ${id} not found`);
  }
  return beacon;
};

export const fetchBeaconFiles = async (
  beaconId: string,
): Promise<BeaconFile[]> => {
  const response = await axios.get<BeaconFile[]>(
    `/api/beaconFiles?BeaconID=${beaconId}`,
  );
  return response.data;
};
