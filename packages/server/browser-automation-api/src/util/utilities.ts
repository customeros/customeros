interface MessageDate {
    year: number;
    month?: number;
    day?: number;
    time: string;
}

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

    static parseDateHeading(heading: string, currentYear: number | null): MessageDate | null {
        const thisYear = new Date().getFullYear();  // Get current year

        // Full date format: "Jul 23, 2023"
        const fullDateMatch = heading.match(/([A-Za-z]+)\s+(\d{1,2}),\s+(\d{4})/);
        if (fullDateMatch) {
            return {
                year: parseInt(fullDateMatch[3]),
                month: this.getMonthNumber(fullDateMatch[1]),
                day: parseInt(fullDateMatch[2]),
                time: ''
            };
        }

        // Month and day format: "Mar 22"
        const monthDayMatch = heading.match(/([A-Za-z]+)\s+(\d{1,2})/);
        if (monthDayMatch) {
            return {
                year: currentYear || thisYear, // Use currentYear if available, otherwise use current year
                month: this.getMonthNumber(monthDayMatch[1]),
                day: parseInt(monthDayMatch[2]),
                time: ''
            };
        }

        // Day of week format: "Wednesday"
        if (heading.match(/^(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday)$/)) {
            const today = new Date();
            const dayOfWeek = heading.toLowerCase();
            const daysOfWeek = ['sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday'];
            const targetDay = daysOfWeek.indexOf(dayOfWeek);
            const currentDay = today.getDay();
            let daysAgo = currentDay - targetDay;
            if (daysAgo <= 0) daysAgo += 7;

            const targetDate = new Date(today);
            targetDate.setDate(today.getDate() - daysAgo);

            return {
                year: targetDate.getFullYear(),
                month: targetDate.getMonth() + 1,
                day: targetDate.getDate(),
                time: ''
            };
        }

        return null;
    }

    static getMonthNumber(monthStr: string): number {
        const months: Record<string, number> = {
            'jan': 1, 'feb': 2, 'mar': 3, 'apr': 4, 'may': 5, 'jun': 6,
            'jul': 7, 'aug': 8, 'sep': 9, 'oct': 10, 'nov': 11, 'dec': 12
        };
        return months[monthStr.toLowerCase()];
    }

    static createDate(messageDate: MessageDate): Date {
        const date = new Date();
        date.setFullYear(messageDate.year);
        if (messageDate.month) date.setMonth(messageDate.month - 1);
        if (messageDate.day) date.setDate(messageDate.day);
        return date;
    }
}