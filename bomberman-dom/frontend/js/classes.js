class Player {
    constructor(id, name, Xposition, Yposition, cellx, celly) {
        this.id = id;
        this.name = name;
        this.Xposition = Xposition;
        this.Yposition = Yposition;
        this.cellx = cellx;
        this.celly = celly;
    }
    move(direction) {
        switch (direction) {
            case "up":
                this.Yposition += 50;
                this.cellx -= 1;
                break;
            case "down":
                this.Yposition -= 50;
                this.cellx += 1;
                break;
            case "left":
                this.Xposition -= 50;
                this.celly -= 1;
                break;
            case "right":
                this.Xposition += 50;
                this.celly += 1;
                break;
        }
    }
}