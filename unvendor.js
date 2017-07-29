/*
    Unvendors by replacing instances of github.com/Necroforger/discordgo with github.com/bwmarrin/discordgo
*/
const fs = require('fs');

function unvendor(path, outdir) {
    fs.readdir(path, (err, files) => {

        if ( err ) {
            console.log(err);
            return;
        }

        fs.mkdir(outdir, () => {
            files.forEach( (file) => {
                fs.readFile(path + "/" + file, "utf8", (err, data) => {
                    fs.writeFile(outdir + "/" + file, data.toString().replace("github.com/Necroforger/discordgo", "github.com/bwmarrin/discordgo"));
                });
            });
        });

    });
}


unvendor("vendor/github.com/Necroforger/dgwidgets", "out");


