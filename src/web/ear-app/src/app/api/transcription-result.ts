export interface TranscriptionResult {
    id: string;
    status: string;
    error: string;
    recognizedText: string;
    progress: number;
}
