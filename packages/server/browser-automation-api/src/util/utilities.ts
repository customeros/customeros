export class TimeUtils {
    static convertToZuluTime(timeStr: string, currentDate: Date | null): string {
        if (!currentDate || !timeStr) return timeStr;

        // Parse time string (e.g., "3:07 PM")
        const [time, period] = timeStr.split(' ');
        const [hours, minutes] = time.split(':').map(num => parseInt(num));

        // Convert to 24-hour format
        let hour24 = hours;
        if (period === 'PM' && hours !== 12) hour24 += 12;
        if (period === 'AM' && hours === 12) hour24 = 0;

        // Create new date with the time
        const dateWithTime = new Date(currentDate);
        dateWithTime.setHours(hour24, minutes, 0, 0);

        // Convert to ISO string (Zulu time)
        return dateWithTime.toISOString();
    }
}