/** @jsxImportSource react */
import Typist from "react-typist";

export default function CodeCard() {
  return (
    <div
      className="mx-auto w-full overflow-hidden rounded-lg sm:w-[600px]"
      aria-hidden="true"
    >
      <div
        className="inverse-toggle h-[300px] overflow-hidden rounded-lg border border-cos-green-200/20 bg-white/10 px-1 pb-6 pt-4 
         font-mono text-[10px] leading-normal text-cos-green-50 subpixel-antialiased shadow-lg transition-all sm:h-[400px] sm:px-2 sm:text-xs md:px-5"
      >
        <div className="top mb-2 flex">
          <div className="h-3 w-3 rounded-full bg-red-500"></div>
          <div className="ml-2 h-3 w-3 rounded-full bg-orange-300"></div>
          <div className="ml-2 h-3 w-3 rounded-full bg-green-500"></div>
        </div>
        <Typist cursor={{ hideWhenDone: true, hideWhenDoneDelay: 0 }}>
        curl http://openline.sh/install.sh | sh
          <Typist.Delay ms={1250} />
        </Typist>
        <Typist
          className="leading-1 translate-y-[-0.2rem] bg-gradient-to-r from-blue-400 via-green-300 to-pink-600 bg-clip-text font-mono text-[7px] text-transparent sm:text-sm md:translate-y-[-0.4rem]"
          cursor={{ show: false }}
          avgTypingDelay={-500}
        >
          <Typist.Delay ms={3500} />
          &nbsp;&nbsp;&nbsp;______&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;__&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;____&nbsp;&nbsp;_____
          <br />
          &nbsp;&nbsp;/&nbsp;____/_&nbsp;&nbsp;_______/&nbsp;/_____&nbsp;&nbsp;____&nbsp;___&nbsp;&nbsp;___&nbsp;&nbsp;_____/&nbsp;__&nbsp;\/&nbsp;___/
          <br />
          &nbsp;/&nbsp;/&nbsp;&nbsp;&nbsp;/&nbsp;/&nbsp;/&nbsp;/&nbsp;___/&nbsp;__/&nbsp;__&nbsp;\/&nbsp;__&nbsp;`__&nbsp;\/&nbsp;_&nbsp;\/&nbsp;___/&nbsp;/&nbsp;/&nbsp;/\__&nbsp;\&nbsp;
          <br />
          /&nbsp;/___/&nbsp;/_/&nbsp;(__&nbsp;&nbsp;)&nbsp;/_/&nbsp;/_/&nbsp;/&nbsp;/&nbsp;/&nbsp;/&nbsp;/&nbsp;/&nbsp;&nbsp;__/&nbsp;/&nbsp;&nbsp;/&nbsp;/_/&nbsp;/___/&nbsp;/&nbsp;
          <br />
          \____/\__,_/____/\__/\____/_/&nbsp;/_/&nbsp;/_/\___/_/&nbsp;&nbsp;&nbsp;\____//____/
          <br />
        </Typist>
        <Typist
          startDelay={3600}
          className=""
          cursor={{ show: false }}
          avgTypingDelay={-100}
        >
          <div>
          VERSION <br />
  openline/0.5.0 darwin-arm64 node-v16.16.0 <br />
  <br />
USAGE <br />
  $ openline [COMMAND] <br />
  <br />
TOPICS <br />
  dev&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;starts and stops local development server for openline applications <br />
  plugins&nbsp;&nbsp;List installed plugins. <br />
  repo&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Get the GitHub repo for an Openline project <br />
  <br />
COMMANDS <br />
  dev&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;starts and stops local development server for openline applications <br />
  help&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Display help for openline. <br />
  issues&nbsp;&nbsp;&nbsp;Interact with GitHub issues for openline-ai.  If no flags are set, <br />
           command will return a list of all open issues assigned to you with a <br />
           milestone or bug tag. <br />
  plugins&nbsp;&nbsp;List installed plugins. <br />
  repo&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Get the GitHub repo for an Openline project <br />
  update&nbsp;&nbsp;&nbsp;update the openline CLI <br />
            <Typist.Delay ms={500} />
          </div>
          <br />
        </Typist>
      </div>
    </div>
  );
}
