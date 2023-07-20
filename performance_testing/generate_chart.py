import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
import mplcursors  # Import the mplcursors library

# Load data from CSV file
data = pd.read_csv('average_response_times.csv')

# Extract the last 10 rows from the data
last_10_rows = data.iloc[-10:]

# Extract x-axis values from the first column
x = last_10_rows.iloc[:, 0]

# Extract y-axis values from all subsequent columns
y = last_10_rows.iloc[:, 1:]

# Set the figure size and margins
fig, ax = plt.subplots(figsize=(16, 12))
fig.subplots_adjust(bottom=0.2)  # Adjust the bottom margin

# Generate the chart
lines = []  # Store the line objects to be used later
for i, column in enumerate(y.columns):
    x_vals = np.arange(len(y[column]))
    line, = plt.plot(x_vals, y[column], marker='o', label=column)
    lines.append(line)  # Append the line object to the list

# Set x-axis tick labels to the values from the first column
plt.xticks(np.arange(len(last_10_rows)), last_10_rows.iloc[:, 0], rotation='vertical')

# Add legend to show the lines for each column
plt.legend(bbox_to_anchor=(-0.3, 1), loc='upper left', borderaxespad=0.)

# Add measurement units to the x-axis and y-axis labels
plt.xlabel('Response time (ms)')
plt.ylabel('Request')

# Save the chart as an image
plt.savefig('average_response_times.png', bbox_inches='tight')

# Use mplcursors to show line names on hover
cursor = mplcursors.cursor(hover=True)
cursor.connect("add", lambda sel: sel.annotation.set_text(lines[sel.index].get_label()))

# Show the plot
plt.show()
