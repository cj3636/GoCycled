#!/bin/bash
# GoCycled Usage Examples
# This script demonstrates various features of the rc command

echo "=== GoCycled (rc) Usage Examples ==="
echo ""

# Check version
echo "1. Check version:"
echo "   $ rc version"
rc version
echo ""

# Show help
echo "2. Show help:"
echo "   $ rc help"
rc help | head -20
echo ""

# Create test files
echo "3. Creating test files..."
mkdir -p /tmp/rc_demo
cd /tmp/rc_demo
echo "Important data" > document.txt
echo "Old log" > old.log
echo "Temporary file" > temp.txt

# Move files to trash
echo "4. Move files to trash:"
echo "   $ rc put document.txt old.log"
rc put document.txt old.log
echo ""

# List trashed items
echo "5. List items in trash:"
echo "   $ rc list"
rc list
echo ""

# Check trash size
echo "6. Check trash size:"
echo "   $ rc size"
rc size
echo ""

# View configuration
echo "7. View configuration:"
echo "   $ rc config"
rc config
echo ""

# Get specific config value
echo "8. Get specific config value:"
echo "   $ rc config get trash_dir"
rc config get trash_dir
echo ""

# Set config value
echo "9. Set config value:"
echo "   $ rc config set confirm_delete false"
rc config set confirm_delete false
echo ""

# Restore a file
echo "10. Restore a file:"
echo "    $ rc restore /tmp/rc_demo/document.txt"
rc restore /tmp/rc_demo/document.txt
echo ""
echo "    Checking if file exists:"
ls -la document.txt 2>/dev/null && echo "    ✓ File restored successfully"
echo ""

# Remove specific item from trash
echo "11. Permanently delete specific item from trash:"
echo "    $ rc remove /tmp/rc_demo/old.log"
rc remove /tmp/rc_demo/old.log
echo ""

# Trash another file for empty demo
echo "12. Trash another file:"
echo "    $ rc put temp.txt"
rc put temp.txt
echo ""

# List trash again
echo "13. List trash again:"
echo "    $ rc list"
rc list
echo ""

# Empty trash (without confirmation since we disabled it)
echo "14. Empty trash:"
echo "    $ rc empty"
rc empty
echo ""

# Verify trash is empty
echo "15. Verify trash is empty:"
echo "    $ rc list"
rc list
echo ""

# Reset config to defaults
echo "16. Reset configuration:"
echo "    $ rc config reset"
rc config reset
echo ""

# Cleanup
cd /
rm -rf /tmp/rc_demo

echo "=== Examples Complete ==="
echo ""
echo "For more information, see the README or run 'rc help'"
