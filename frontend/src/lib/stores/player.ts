import { writable } from "svelte/store";

interface AudioFile {
  title: string;
  artist: string;
  album: string;
  album_artist: string;
  genre: string;
  year: number;
  track_num: number;
  track_total: number;
  disc_num: number;
  disc_total: number;
  has_art: boolean;
  file_path: string;
  ext: string;
}

interface PlayerState {
  track: AudioFile | null;
  queue: AudioFile[];
  queueIndex: number;
  paused: boolean;
  currentTime: number;
  duration: number;
  volume: number;
}

const initial: PlayerState = {
  track: null,
  queue: [],
  queueIndex: -1,
  paused: true,
  currentTime: 0,
  duration: 0,
  volume: parseFloat(localStorage.getItem("harmonic:volume") ?? "0.8"),
};

export const player = writable<PlayerState>(initial);

export function playTrack(track: AudioFile, queue: AudioFile[], index: number) {
  player.update(s => ({ ...s, track, queue, queueIndex: index, paused: false, currentTime: 0, duration: 0 }));
}


export function togglePause() {
  player.update(s => s.track ? { ...s, paused: !s.paused } : s);
}

export function nextTrack() {
  player.update(s => {
    if (s.queueIndex < s.queue.length - 1) {
      const i = s.queueIndex + 1;
      return { ...s, queueIndex: i, track: s.queue[i], currentTime: 0, duration: 0 };
    }
    return { ...s, track: null, paused: true, currentTime: 0, duration: 0, queueIndex: -1 };
  });
}

export function prevTrack() {
  player.update(s => {
    if (s.queueIndex > 0) {
      const i = s.queueIndex - 1;
      return { ...s, queueIndex: i, track: s.queue[i], currentTime: 0, duration: 0 };
    }
    return s;
  });
}

export function setVolume(vol: number) {
  player.update(s => ({ ...s, volume: Math.max(0, Math.min(1, vol)) }));
  localStorage.setItem("harmonic:volume", Math.max(0, Math.min(1, vol)).toString());
}

let _preMuteVolume = 0.8;
export function toggleMute() {
  player.update(s => {
    if (s.volume > 0) {
      _preMuteVolume = s.volume;
      return { ...s, volume: 0 };
    }
    return { ...s, volume: _preMuteVolume };
  });
}

export function audioUrl(track: AudioFile): string {
  return "/audio?path=" + encodeURIComponent(track.file_path);
}